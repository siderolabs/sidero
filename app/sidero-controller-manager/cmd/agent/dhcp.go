// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/nclient4"
	"github.com/jsimonetti/rtnetlink"
	"github.com/siderolabs/gen/xslices"
	"github.com/siderolabs/go-procfs/procfs"
	"github.com/siderolabs/go-retry/retry"
	"golang.org/x/sys/unix"

	"github.com/siderolabs/sidero/app/sidero-controller-manager/pkg/constants"
)

func setupNetworking() error {
	var bootMACAddress string

	if found := procfs.ProcCmdline().Get(constants.AgentMACArg).First(); found != nil {
		bootMACAddress = *found
	} else {
		return fmt.Errorf("no MAC found")
	}

	link, err := waitForLink(bootMACAddress)
	if err != nil {
		return err
	}

	if err = brinkLinkUp(link.Index); err != nil {
		return err
	}

	return runDHCP(link)
}

func waitForLink(hwaddr string) (net.Interface, error) {
	log.Printf("waiting for network link with MAC %q...", hwaddr)

	var foundLink net.Interface

	err := retry.Constant(time.Minute, retry.WithUnits(time.Second)).Retry(func() error {
		links, err := net.Interfaces()
		if err != nil {
			return err
		}

		for _, link := range links {
			if link.HardwareAddr.String() == hwaddr {
				foundLink = link

				return nil
			}
		}

		return retry.ExpectedErrorf("link with MAC %q not found", hwaddr)
	})

	return foundLink, err
}

func brinkLinkUp(linkIndex int) error {
	log.Printf("bringing link up...")

	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		return fmt.Errorf("error dialing rtnetlink socket: %w", err)
	}

	defer conn.Close() //nolint:errcheck

	if err := conn.Link.Set(&rtnetlink.LinkMessage{
		Family: unix.AF_UNSPEC,
		Type:   unix.ARPHRD_ETHER,
		Index:  uint32(linkIndex),
		Flags:  unix.IFF_UP,
		Change: unix.IFF_UP,
	}); err != nil {
		return fmt.Errorf("error setting link up: %w", err)
	}

	return nil
}

func runDHCP(link net.Interface) error {
	log.Printf("running DHCP on %q...", link.Name)

	var lease *nclient4.Lease

	if err := retry.Constant(5*time.Minute,
		retry.WithUnits(10*time.Second),
		retry.WithAttemptTimeout(30*time.Second),
	).RetryWithContext(context.Background(), func(ctx context.Context) error {
		var err error

		lease, err = acquireLease(ctx, link.Name)
		if err != nil {
			return retry.ExpectedError(err)
		}

		return nil
	}); err != nil {
		return err
	}

	log.Printf("got DHCP lease: %s", lease.ACK.Summary())

	return configureNetworking(link.Index, lease)
}

func acquireLease(ctx context.Context, linkName string) (*nclient4.Lease, error) {
	opts := []dhcpv4.OptionCode{
		dhcpv4.OptionClasslessStaticRoute,
		dhcpv4.OptionDomainNameServer,
		dhcpv4.OptionInterfaceMTU,
		dhcpv4.OptionHostName,
	}

	mods := []dhcpv4.Modifier{dhcpv4.WithRequestedOptions(opts...)}

	cli, err := nclient4.New(linkName)
	if err != nil {
		return nil, fmt.Errorf("error creating DHCP client: %w", err)
	}

	//nolint:errcheck
	defer cli.Close()

	lease, err := cli.Request(ctx, mods...)
	if err != nil {
		return nil, fmt.Errorf("error requesting DHCP lease: %w", err)
	}

	return lease, nil
}

func configureNetworking(linkIndex int, lease *nclient4.Lease) error {
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		return fmt.Errorf("error dialing rtnetlink socket: %w", err)
	}

	defer conn.Close() //nolint:errcheck

	if err := conn.Link.Set(&rtnetlink.LinkMessage{
		Family: unix.AF_UNSPEC,
		Type:   unix.ARPHRD_ETHER,
		Index:  uint32(linkIndex),
		Flags:  unix.IFF_UP,
		Change: unix.IFF_UP,
	}); err != nil {
		return fmt.Errorf("error setting link up: %w", err)
	}

	prefixLen, _ := lease.ACK.SubnetMask().Size()

	log.Printf("assigning address %s/%d", lease.ACK.YourIPAddr, prefixLen)

	if err := conn.Address.New(&rtnetlink.AddressMessage{
		Family:       unix.AF_INET,
		PrefixLength: uint8(prefixLen),
		Scope:        unix.RT_SCOPE_UNIVERSE,
		Index:        uint32(linkIndex),
		Attributes: &rtnetlink.AddressAttributes{
			Address: lease.ACK.YourIPAddr,
			Local:   lease.ACK.YourIPAddr,
		},
	}); err != nil {
		return fmt.Errorf("error adding address: %w", err)
	}

	mtu, err := dhcpv4.GetUint16(dhcpv4.OptionInterfaceMTU, lease.ACK.Options)
	if err == nil {
		log.Printf("setting MTU to %d", mtu)

		if err := conn.Link.Set(&rtnetlink.LinkMessage{
			Family: unix.AF_UNSPEC,
			Type:   unix.ARPHRD_ETHER,
			Index:  uint32(linkIndex),
			Attributes: &rtnetlink.LinkAttributes{
				MTU: uint32(mtu),
			},
		}); err != nil {
			return fmt.Errorf("error setting MTU: %w", err)
		}
	}

	addRoute := func(destination *net.IPNet, gateway net.IP) error {
		log.Printf("adding route %s via %s", destination, gateway)

		var dstLength int

		if destination != nil {
			dstLength, _ = destination.Mask.Size()
		}

		var dstAddr net.IP

		if destination != nil {
			dstAddr = destination.IP
		}

		if err := conn.Route.Add(&rtnetlink.RouteMessage{
			Family:    unix.AF_INET,
			DstLength: uint8(dstLength),
			Scope:     unix.RT_SCOPE_UNIVERSE,
			Table:     unix.RT_TABLE_MAIN,
			Protocol:  unix.RTPROT_BOOT,
			Type:      unix.RTN_UNICAST,
			Attributes: rtnetlink.RouteAttributes{
				Dst:      dstAddr,
				Gateway:  gateway,
				OutIface: uint32(linkIndex),
			},
		}); err != nil {
			return fmt.Errorf("error adding route: %w", err)
		}

		return nil
	}

	if len(lease.ACK.ClasslessStaticRoute()) > 0 {
		for _, route := range lease.ACK.ClasslessStaticRoute() {
			if err := addRoute(route.Dest, route.Router); err != nil {
				return err
			}
		}
	} else {
		for _, router := range lease.ACK.Router() {
			if err := addRoute(nil, router); err != nil {
				return err
			}
		}
	}

	if lease.ACK.HostName() != "" {
		log.Printf("setting hostname to %q", lease.ACK.HostName())

		unix.Sethostname([]byte(lease.ACK.HostName())) //nolint:errcheck
	}

	if len(lease.ACK.DNS()) > 0 {
		log.Printf("setting DNS servers to %s", lease.ACK.DNS())

		contents := strings.Join(xslices.Map(lease.ACK.DNS(),
			func(ns net.IP) string {
				return fmt.Sprintf("nameserver %s\n", ns)
			}), "")

		if err := os.WriteFile("/etc/resolv.conf", []byte(contents), 0o777); err != nil {
			return fmt.Errorf("error writing /etc/resolv.conf: %w", err)
		}
	}

	return nil
}
