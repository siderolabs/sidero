// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"log"
	"net"

	"github.com/siderolabs/go-blockdevice/blockdevice/util/disk"
	"github.com/siderolabs/go-smbios/smbios"

	"github.com/siderolabs/sidero/app/sidero-controller-manager/internal/api"
)

func MapHardwareInformation(s *smbios.SMBIOS, disks []*disk.Disk, interfaces []net.Interface) *api.HardwareInformation {
	if s != nil {
		return &api.HardwareInformation{
			System:  MapSystemInformation(s),
			Compute: MapComputeInformation(s),
			Memory:  MapMemoryInformation(s),
			Storage: MapStorageInformation(disks),
			Network: MapNetworkInformation(interfaces),
		}
	}

	return &api.HardwareInformation{
		Storage: MapStorageInformation(disks),
	}
}

func MapSystemInformation(s *smbios.SMBIOS) *api.SystemInformation {
	return &api.SystemInformation{
		Manufacturer: s.SystemInformation.Manufacturer,
		ProductName:  s.SystemInformation.ProductName,
		SerialNumber: s.SystemInformation.SerialNumber,
		Uuid:         s.SystemInformation.UUID,
		SkuNumber:    s.SystemInformation.SKUNumber,
		Family:       s.SystemInformation.Family,
	}
}

func MapComputeInformation(s *smbios.SMBIOS) *api.ComputeInformation {
	var (
		totalCoreCount   = 0
		totalThreadCount = 0
		processors       []*api.Processor
	)

	for _, v := range s.ProcessorInformation {
		if v.Status.SocketPopulated() {
			totalCoreCount += int(v.CoreCount)
			totalThreadCount += int(v.ThreadCount)

			processor := &api.Processor{
				Manufacturer: v.ProcessorManufacturer,
				ProductName:  v.ProcessorVersion,
				SerialNumber: v.SerialNumber,
				Speed:        uint32(v.CurrentSpeed),
				CoreCount:    uint32(v.CoreCount),
				ThreadCount:  uint32(v.ThreadCount),
			}

			processors = append(processors, processor)
		}
	}

	return &api.ComputeInformation{
		TotalCoreCount:   uint32(totalCoreCount),
		TotalThreadCount: uint32(totalThreadCount),
		ProcessorCount:   uint32(len(processors)),
		Processors:       processors,
	}
}

func MapMemoryInformation(s *smbios.SMBIOS) *api.MemoryInformation {
	var (
		totalSize = 0
		modules   []*api.MemoryModule
	)

	for _, v := range s.MemoryDevices {
		if v.Size != 0 && v.Size != 0xFFFF {
			var size uint32

			if v.Size == 0x7FFF {
				totalSize += int(v.ExtendedSize)
				size = uint32(v.ExtendedSize)
			} else {
				totalSize += v.Size.Megabytes()
				size = uint32(v.Size)
			}

			memoryModule := &api.MemoryModule{
				Manufacturer: v.Manufacturer,
				ProductName:  v.PartNumber,
				SerialNumber: v.SerialNumber,
				Type:         v.MemoryType.String(),
				Size:         size,
				Speed:        uint32(v.Speed),
			}

			modules = append(modules, memoryModule)
		}
	}

	return &api.MemoryInformation{
		TotalSize:   uint32(totalSize),
		ModuleCount: uint32(len(modules)),
		Modules:     modules,
	}
}

func MapStorageInformation(s []*disk.Disk) *api.StorageInformation {
	totalSize := uint64(0)
	devices := make([]*api.StorageDevice, 0, len(s))

	for _, v := range s {
		totalSize += v.Size

		var storageType api.StorageType

		switch v.Type.String() {
		case "ssd":
			storageType = api.StorageType_SSD
		case "hdd":
			storageType = api.StorageType_HDD
		case "nvme":
			storageType = api.StorageType_NVMe
		case "sd":
			storageType = api.StorageType_SD
		default:
			storageType = api.StorageType_Unknown
		}

		storageDevice := &api.StorageDevice{
			Model:      v.Model,
			Serial:     v.Serial,
			Type:       storageType,
			Size:       v.Size,
			Name:       v.Name,
			DeviceName: v.DeviceName,
			Uuid:       v.UUID,
			Wwid:       v.WWID,
		}

		devices = append(devices, storageDevice)
	}

	return &api.StorageInformation{
		TotalSize:   totalSize,
		DeviceCount: uint32(len(devices)),
		Devices:     devices,
	}
}

func MapNetworkInformation(s []net.Interface) *api.NetworkInformation {
	interfaces := make([]*api.NetworkInterface, 0, len(s))

	for _, v := range s {
		if len(v.HardwareAddr) == 0 {
			continue // skip interfaces without a hardware address
		}

		addrs, err := v.Addrs()
		if err != nil {
			log.Printf("encountered error fetching addresses of network interface %q: %q", v.Name, err)

			addrs = make([]net.Addr, 0)
		}

		var addresses []string

		for _, a := range addrs {
			addresses = append(addresses, a.String())
		}

		networkInterface := &api.NetworkInterface{
			Index:     uint32(v.Index),
			Name:      v.Name,
			Flags:     v.Flags.String(),
			Mtu:       uint32(v.MTU),
			Mac:       v.HardwareAddr.String(),
			Addresses: addresses,
		}

		interfaces = append(interfaces, networkInterface)
	}

	return &api.NetworkInformation{
		InterfaceCount: uint32(len(interfaces)),
		Interfaces:     interfaces,
	}
}
