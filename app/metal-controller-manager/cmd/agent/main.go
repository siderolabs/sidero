// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/talos-systems/go-blockdevice/blockdevice/probe"
	"github.com/talos-systems/go-procfs/procfs"
	"github.com/talos-systems/go-smbios/smbios"
	talosnet "github.com/talos-systems/net"
	"golang.org/x/sys/unix"
	"google.golang.org/grpc"

	"github.com/talos-systems/sidero/app/metal-controller-manager/internal/api"
	"github.com/talos-systems/sidero/app/metal-controller-manager/pkg/constants"
)

func setup() error {
	if err := os.MkdirAll("/dev", 0o777); err != nil {
		return err
	}

	if err := os.MkdirAll("/proc", 0o777); err != nil {
		return err
	}

	if err := os.MkdirAll("/sys", 0o777); err != nil {
		return err
	}

	if err := os.MkdirAll("/tmp", 0o777); err != nil {
		return err
	}

	if err := unix.Mount("devtmpfs", "/dev", "devtmpfs", unix.MS_NOSUID, "mode=0755"); err != nil {
		return err
	}

	if err := unix.Mount("proc", "/proc", "proc", unix.MS_NOSUID|unix.MS_NOEXEC|unix.MS_NODEV, ""); err != nil {
		return err
	}

	if err := unix.Mount("sysfs", "/sys", "sysfs", 0, ""); err != nil {
		return err
	}

	if err := unix.Mount("tmpfs", "/tmp", "tmpfs", 0, ""); err != nil {
		return err
	}

	kmsg, err := os.OpenFile("/dev/kmsg", os.O_RDWR|unix.O_CLOEXEC|unix.O_NONBLOCK|unix.O_NOCTTY, 0o666)
	if err != nil {
		return fmt.Errorf("failed to open /dev/kmsg: %w", err)
	}

	log.SetOutput(kmsg)
	log.SetPrefix("[sidero]" + " ")
	log.SetFlags(0)

	return nil
}

func create(endpoint string, s *smbios.Smbios) (*api.CreateServerResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, endpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	c := api.NewAgentClient(conn)

	uuid, err := s.SystemInformation().UUID()
	if err != nil {
		return nil, err
	}

	req := &api.CreateServerRequest{
		SystemInformation: &api.SystemInformation{
			Uuid:         uuid.String(),
			Manufacturer: s.SystemInformation().Manufacturer(),
			ProductName:  s.SystemInformation().ProductName(),
			Version:      s.SystemInformation().Version(),
			SerialNumber: s.SystemInformation().SerialNumber(),
			SkuNumber:    s.SystemInformation().SKUNumber(),
			Family:       s.SystemInformation().Family(),
		},
		Cpu: &api.CPU{
			Manufacturer: s.ProcessorInformation().ProcessorManufacturer(),
			Version:      s.ProcessorInformation().ProcessorVersion(),
		},
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("encountered error fetching hostname: %q", err)
	} else {
		req.Hostname = hostname
	}

	resp, err := c.CreateServer(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func wipe(endpoint string, s *smbios.Smbios) error {
	uuid, err := s.SystemInformation().UUID()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, endpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()

	c := api.NewAgentClient(conn)

	_, err = c.MarkServerAsWiped(ctx, &api.MarkServerAsWipedRequest{Uuid: uuid.String()})
	if err != nil {
		return err
	}

	return nil
}

func reconcileIPs(endpoint string, s *smbios.Smbios, ips []net.IP) error {
	uuid, err := s.SystemInformation().UUID()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, endpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()

	c := api.NewAgentClient(conn)

	addresses := make([]*api.Address, len(ips))
	for i := range addresses {
		addresses[i] = &api.Address{
			Type:    "InternalIP",
			Address: ips[i].String(),
		}
	}

	_, err = c.ReconcileServerAddresses(ctx, &api.ReconcileServerAddressesRequest{
		Uuid:    uuid.String(),
		Address: addresses,
	})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := setup(); err != nil {
		log.Fatal(err)
	}

	var endpoint *string
	if endpoint = procfs.ProcCmdline().Get(constants.AgentEndpointArg).First(); endpoint == nil {
		log.Fatal(fmt.Errorf("no endpoint found"))
	}

	log.Printf("Using %q as API endpoint", *endpoint)

	log.Println("Reading SMBIOS")

	s, err := smbios.New()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := create(*endpoint, s)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Registration complete")

	if resp.GetWipe() {
		bds, err := probe.All()
		if err != nil {
			log.Fatal(err)
		}

		for _, bd := range bds {
			log.Printf("Resetting %s", bd.Path)

			_, err := bd.Device().WriteAt(bytes.Repeat([]byte{0}, 512), 0)
			if err != nil {
				log.Fatal(err)
			}

			if err := bd.Reset(); err != nil {
				log.Fatal(err)
			}

			if err := bd.Close(); err != nil {
				log.Fatal(err)
			}
		}

		if err := wipe(*endpoint, s); err != nil {
			log.Fatal(err)
		}

		log.Println("Wipe complete")
	}

	ips, err := talosnet.IPAddrs()
	if err != nil {
		log.Println("failed to discover IPs")
	} else {
		if err := reconcileIPs(*endpoint, s, ips); err != nil {
			log.Fatal(err)
		}

		log.Printf("Reconciled IPs")
	}

	// nolint: errcheck
	unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF)
}
