// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

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

	BootstrapProviders      []string
	InfrastructureProviders []string
	ControlPlaneProviders   []string

	RegistryMirrors stringSlice

	ManagementCIDR  string
	ManagementNodes int

	MemMB  int64
	CPUs   int64
	DiskGB int64

	TalosctlPath string
}

const defaulTalosRelease = "v0.7.0-alpha.2"

// DefaultOptions returns default settings.
func DefaultOptions() Options {
	return Options{
		BootstrapClusterName:    "sfyra",
		BootstrapTalosVmlinuz:   "_out/vmlinuz",
		BootstrapTalosInitramfs: "_out/initramfs.xz",
		BootstrapTalosInstaller: fmt.Sprintf("docker.io/autonomy/installer:%s", defaulTalosRelease),
		BootstrapCIDR:           "172.24.0.0/24",

		TalosKernelURL: fmt.Sprintf("https://github.com/talos-systems/talos/releases/download/%s/vmlinuz", defaulTalosRelease),
		TalosInitrdURL: fmt.Sprintf("https://github.com/talos-systems/talos/releases/download/%s/initramfs.xz", defaulTalosRelease),
		TalosInstaller: fmt.Sprintf("docker.io/autonomy/installer:%s", defaulTalosRelease),

		BootstrapProviders:      []string{"talos"},
		InfrastructureProviders: []string{"sidero"},
		ControlPlaneProviders:   []string{"talos"},

		ManagementCIDR:  "172.25.0.0/24",
		ManagementNodes: 4,

		MemMB:  2048,
		CPUs:   2,
		DiskGB: 4,

		TalosctlPath: "_out/talosctl-linux-amd64",
	}
}
