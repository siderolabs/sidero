/*
Copyright 2020 Talos Systems, Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/talos-systems/go-procfs/procfs"
	"github.com/talos-systems/go-smbios/smbios"
	"github.com/talos-systems/sidero/internal/app/metal-controller-manager/internal/api"
	"golang.org/x/sys/unix"
	"google.golang.org/grpc"
)

func print(s *smbios.Smbios) {
	log.Println(s.BIOSInformation().Vendor())
	log.Println(s.BIOSInformation().Version())
	log.Println(s.BIOSInformation().ReleaseDate())
	log.Println(s.SystemInformation().Manufacturer())
	log.Println(s.SystemInformation().ProductName())
	log.Println(s.SystemInformation().Version())
	log.Println(s.SystemInformation().SerialNumber())
	log.Println(s.SystemInformation().SKUNumber())
	log.Println(s.SystemInformation().Family())
	log.Println(s.BaseboardInformation().Manufacturer())
	log.Println(s.BaseboardInformation().Product())
	log.Println(s.BaseboardInformation().Version())
	log.Println(s.BaseboardInformation().SerialNumber())
	log.Println(s.BaseboardInformation().AssetTag())
	log.Println(s.BaseboardInformation().LocationInChassis())
	log.Println(s.SystemEnclosure().Manufacturer())
	log.Println(s.SystemEnclosure().Version())
	log.Println(s.SystemEnclosure().SerialNumber())
	log.Println(s.SystemEnclosure().AssetTagNumber())
	log.Println(s.SystemEnclosure().SKUNumber())
	log.Println(s.ProcessorInformation().SocketDesignation())
	log.Println(s.ProcessorInformation().ProcessorManufacturer())
	log.Println(s.ProcessorInformation().ProcessorVersion())
	log.Println(s.ProcessorInformation().SerialNumber())
	log.Println(s.ProcessorInformation().AssetTag())
	log.Println(s.ProcessorInformation().PartNumber())
	log.Println(s.CacheInformation().SocketDesignation())
	log.Println(s.PortConnectorInformation().InternalReferenceDesignator())
	log.Println(s.PortConnectorInformation().ExternalReferenceDesignator())
	log.Println(s.SystemSlots().SlotDesignation())
	log.Println(s.BIOSLanguageInformation().CurrentLanguage())
	log.Println(s.GroupAssociations().GroupName())
}

func create(endpoint string, s *smbios.Smbios) error {
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()

	c := api.NewDiscoveryClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

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
	if err := os.MkdirAll("/dev", 0777); err != nil {
		return err
	}

	if err := os.MkdirAll("/proc", 0777); err != nil {
		return err
	}

	if err := os.MkdirAll("/sys", 0777); err != nil {
		return err
	}

	if err := os.MkdirAll("/tmp", 0777); err != nil {
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

	kmsg, err := os.OpenFile("/dev/kmsg", os.O_RDWR|unix.O_CLOEXEC|unix.O_NONBLOCK|unix.O_NOCTTY, 0666)
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

	unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF)
}
