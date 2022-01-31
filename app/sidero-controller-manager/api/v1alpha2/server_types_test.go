// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

//nolint:scopelint
package v1alpha2_test

import (
	"testing"

	metal "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha2"
)

func Test_PartialEqual(t *testing.T) {
	info := &metal.HardwareInformation{
		System: &metal.SystemInformation{
			Uuid:         "4c4c4544-0039-3010-8048-b7c04f384432",
			Manufacturer: "Dell Inc.",
			ProductName:  "PowerEdge R630",
			SerialNumber: "790H8D2",
			SKUNumber:    "",
			Family:       "",
		},
		Compute: &metal.ComputeInformation{
			TotalCoreCount:   8,
			TotalThreadCount: 16,
			ProcessorCount:   1,
			Processors: []*metal.Processor{
				{
					Manufacturer: "Intel",
					ProductName:  "Intel(R) Xeon(R) CPU E5-2630 v3 @ 2.40GHz",
					SerialNumber: "",
					Speed:        2400,
					CoreCount:    8,
					ThreadCount:  16,
				},
			},
		},
		Memory: &metal.MemoryInformation{
			TotalSize:   "32 GB",
			ModuleCount: 2,
			Modules: []*metal.MemoryModule{
				{
					Manufacturer: "002C00B3002C",
					ProductName:  "18ASF2G72PDZ-2G3B1",
					SerialNumber: "12BDC045",
					Type:         "LPDDR3",
					Size:         16384,
					Speed:        2400,
				},
				{
					Manufacturer: "002C00B3002C",
					ProductName:  "18ASF2G72PDZ-2G3B1",
					SerialNumber: "12BDBF5D",
					Type:         "LPDDR3",
					Size:         16384,
					Speed:        2400,
				},
			},
		},
		Storage: &metal.StorageInformation{
			TotalSize:   "1116 GB",
			DeviceCount: 1,
			Devices: []*metal.StorageDevice{
				{
					Type:       "HDD",
					Size:       1199101181952,
					Model:      "PERC H730 Mini",
					Serial:     "",
					Name:       "sda",
					DeviceName: "/dev/sda",
					UUID:       "",
					WWID:       "naa.61866da055de070028d8e83307cc6df2",
				},
			},
		},
		Network: &metal.NetworkInformation{
			InterfaceCount: 2,
			Interfaces: []*metal.NetworkInterface{
				{
					Index:     1,
					Name:      "lo",
					Flags:     "up|loopback",
					MTU:       65536,
					MAC:       "",
					Addresses: []string{"127.0.0.1/8", "::1/128"},
				},
				{
					Index:     2,
					Name:      "enp3s0",
					Flags:     "up|broadcast|multicast",
					MTU:       1500,
					MAC:       "40:8d:5c:86:5a:14",
					Addresses: []string{"192.168.2.8/24", "fe80::dcb3:295c:755b:91bb/64"},
				},
			},
		},
	}

	tests := []struct {
		name string
		args *metal.HardwareInformation
		want bool
	}{
		{
			name: "defaults are partially equal",
			args: &metal.HardwareInformation{},
			want: true,
		},
		{
			name: "cpu is partially equal",
			args: &metal.HardwareInformation{
				Compute: &metal.ComputeInformation{
					Processors: []*metal.Processor{
						{
							Manufacturer: "Intel",
						},
					},
				},
				// Skip all other fields to indicate that we don't want to compare it.
			},
			want: true,
		},
		{
			name: "cpu is not partially equal",
			args: &metal.HardwareInformation{
				Compute: &metal.ComputeInformation{
					Processors: []*metal.Processor{
						{
							Manufacturer: "AMD",
						},
					},
				},
				// Skip all other fields to indicate that we don't want to compare it.
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.PartialEqual(info); got != tt.want {
				t.Errorf("PartialEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
