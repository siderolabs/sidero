// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import "fmt"

// Options control the sidero testing.
type Options struct {
	SkipTeardown bool

	BootstrapClusterName    string
	BootstrapTalosVmlinuz   string
	BootstrapTalosInitramfs string
	BootstrapTalosInstaller string
	BootstrapCIDR           string

	TalosKernelURL string
	TalosInitrdURL string
	TalosInstaller string

	ClusterctlConfigPath    string
	BootstrapProviders      []string
	InfrastructureProviders []string
	ControlPlaneProviders   []string

	RegistryMirrors []string

	ManagementCIDR    string
	ManagementSetName string
	ManagementNodes   int

	MemMB  int64
	CPUs   int64
	DiskGB int64

	TalosctlPath string
}

// TalosRelease is set as build argument.
var (
	TalosRelease string
)

// DefaultOptions returns default settings.
func DefaultOptions() Options {
	return Options{
		BootstrapClusterName:    "sfyra",
		BootstrapTalosVmlinuz:   fmt.Sprintf("_out/%s/vmlinuz", TalosRelease),
		BootstrapTalosInitramfs: fmt.Sprintf("_out/%s/initramfs.xz", TalosRelease),
		BootstrapTalosInstaller: fmt.Sprintf("docker.io/autonomy/installer:%s", TalosRelease),
		BootstrapCIDR:           "172.24.0.0/24",

		TalosKernelURL: fmt.Sprintf("https://github.com/talos-systems/talos/releases/download/%s/vmlinuz", TalosRelease),
		TalosInitrdURL: fmt.Sprintf("https://github.com/talos-systems/talos/releases/download/%s/initramfs.xz", TalosRelease),
		TalosInstaller: fmt.Sprintf("docker.io/autonomy/installer:%s", TalosRelease),

		BootstrapProviders:      []string{"talos"},
		InfrastructureProviders: []string{"sidero"},
		ControlPlaneProviders:   []string{"talos"},

		ManagementCIDR:    "172.25.0.0/24",
		ManagementSetName: "sfyra-management",
		ManagementNodes:   4,

		MemMB:  2048,
		CPUs:   2,
		DiskGB: 4,

		TalosctlPath: fmt.Sprintf("_out/%s/talosctl-linux-amd64", TalosRelease),
	}
}
