// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/talos-systems/go-procfs/procfs"
	"github.com/talos-systems/go-smbios/smbios"
	"golang.org/x/sys/unix"
	"google.golang.org/grpc"

	"github.com/talos-systems/sidero/internal/app/metal-controller-manager/internal/api"
)

func create(endpoint string, s *smbios.Smbios) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, endpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()

	c := api.NewDiscoveryClient(conn)

	uuid, err := s.SystemInformation().UUID()
	if err != nil {
		return err
	}

	_, err = c.CreateServer(ctx, &api.CreateServerRequest{
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
	})
	if err != nil {
		return err
	}

	return nil
}

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

func main() {
	if err := setup(); err != nil {
		log.Fatal(err)
	}

	log.Println("Reading SMBIOS")

	s, err := smbios.New()
	if err != nil {
		log.Fatal(err)
	}

	var endpoint *string
	if endpoint = procfs.ProcCmdline().Get("sidero.endpoint").First(); endpoint == nil {
		log.Fatal(fmt.Errorf("no endpoint found"))
	}

	log.Printf("Creating resource via %q", *endpoint)

	if err = create(*endpoint, s); err != nil {
		log.Fatal(err)
	}

	log.Println("Discovery complete")

	// nolint: errcheck
	unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF)
}
