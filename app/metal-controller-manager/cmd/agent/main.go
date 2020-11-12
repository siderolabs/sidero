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
	"time"

	"github.com/talos-systems/go-blockdevice/blockdevice"
	"github.com/talos-systems/go-blockdevice/blockdevice/util"
	"github.com/talos-systems/go-procfs/procfs"
	"github.com/talos-systems/go-retry/retry"
	"github.com/talos-systems/go-smbios/smbios"
	talosnet "github.com/talos-systems/net"
	"golang.org/x/sync/errgroup"
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

func create(ctx context.Context, client api.AgentClient, s *smbios.Smbios) (*api.CreateServerResponse, error) {
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

	var resp *api.CreateServerResponse

	err = retry.Constant(5*time.Minute, retry.WithUnits(30*time.Second)).Retry(func() error {
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		resp, err = client.CreateServer(ctx, req)
		if err != nil {
			return retry.ExpectedError(err)
		}

		return nil
	})

	return resp, err
}

func wipe(ctx context.Context, client api.AgentClient, s *smbios.Smbios) error {
	uuid, err := s.SystemInformation().UUID()
	if err != nil {
		return err
	}

	return retry.Constant(5*time.Minute, retry.WithUnits(30*time.Second)).Retry(func() error {
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		_, err = client.MarkServerAsWiped(ctx, &api.MarkServerAsWipedRequest{Uuid: uuid.String()})
		if err != nil {
			return retry.ExpectedError(err)
		}

		return nil
	})
}

func reconcileIPs(ctx context.Context, client api.AgentClient, s *smbios.Smbios, ips []net.IP) error {
	uuid, err := s.SystemInformation().UUID()
	if err != nil {
		return err
	}

	addresses := make([]*api.Address, len(ips))
	for i := range addresses {
		addresses[i] = &api.Address{
			Type:    "InternalIP",
			Address: ips[i].String(),
		}
	}

	return retry.Constant(5*time.Minute, retry.WithUnits(30*time.Second)).Retry(func() error {
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		_, err = client.ReconcileServerAddresses(ctx, &api.ReconcileServerAddressesRequest{
			Uuid:    uuid.String(),
			Address: addresses,
		})
		if err != nil {
			return retry.ExpectedError(err)
		}

		return nil
	})
}

func shutdown(err error) {
	if err != nil {
		log.Println(err)
	}

	for i := 10; i >= 0; i-- {
		log.Printf("rebooting in %d seconds\n", i)
		time.Sleep(1 * time.Second)
	}

	if unix.Reboot(unix.LINUX_REBOOT_CMD_RESTART) == nil {
		select {}
	}

	os.Exit(1)
}

func connect(ctx context.Context, endpoint string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return grpc.DialContext(ctx, endpoint, grpc.WithInsecure())
}

func mainFunc() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := setup(); err != nil {
		return err
	}

	var endpoint string
	if found := procfs.ProcCmdline().Get(constants.AgentEndpointArg).First(); found != nil {
		endpoint = *found
	} else {
		return fmt.Errorf("no endpoint found")
	}

	log.Printf("Using %q as API endpoint", endpoint)

	conn, err := connect(ctx, endpoint)
	if err != nil {
		return err
	}

	defer conn.Close()

	client := api.NewAgentClient(conn)

	log.Println("Reading SMBIOS")

	s, err := smbios.New()
	if err != nil {
		return err
	}

	createResp, err := create(ctx, client, s)
	if err != nil {
		return err
	}

	log.Println("Registration complete")

	ips, err := talosnet.IPAddrs()
	if err != nil {
		log.Println("failed to discover IPs")
	} else {
		if err = reconcileIPs(ctx, client, s, ips); err != nil {
			shutdown(err)
		}

		log.Printf("Reconciled IPs")
	}

	if createResp.GetWipe() {
		disks, err := util.GetDisks()
		if err != nil {
			shutdown(err)
		}

		var eg errgroup.Group

		for _, disk := range disks {
			func(path string) {
				eg.Go(func() error {
					log.Printf("Resetting %s", path)

					bd, err := blockdevice.Open(path)
					if err != nil {
						log.Printf("Skipping %s: %s", path, err)

						return nil
					}

					if createResp.GetInsecureWipe() {
						if err = bd.FastWipe(); err != nil {
							return fmt.Errorf("failed wiping %q: %w", path, err)
						}

						log.Printf("Fast wiped %s", path)
					} else {
						method, err := bd.Wipe()
						if err != nil {
							return fmt.Errorf("failed wiping %q: %w", path, err)
						}

						log.Printf("Wiped %s with %s", path, method)
					}

					return bd.Close()
				})
			}(disk.DeviceName)
		}

		if err := eg.Wait(); err != nil {
			shutdown(err)
		}

		if err := wipe(ctx, client, s); err != nil {
			shutdown(err)
		}

		log.Println("Wipe complete")
	}

	return nil
}

func main() {
	shutdown(mainFunc())
}
