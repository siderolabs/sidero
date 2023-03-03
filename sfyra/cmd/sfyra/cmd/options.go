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
	BootstrapCNIBundleURL   string
	BootstrapCIDR           string

	TalosKernelURL string
	TalosInitrdURL string

	ClusterctlConfigPath    string
	CoreProvider            string
	BootstrapProviders      []string
	InfrastructureProviders []string
	ControlPlaneProviders   []string

	RegistryMirrors []string

	ManagementNodes int

	BootstrapMemMB  int64
	BootstrapCPUs   int64
	BootstrapDiskGB int64

	ManagementMemMB  int64
	ManagementCPUs   int64
	ManagementDiskGB int64

	DefaultBootOrder string

	TalosctlPath string

	PowerSimulatedExplicitFailureProb float64
	PowerSimulatedSilentFailureProb   float64
}

// TalosRelease and KubernetesVersion are set as build arguments.
var (
	TalosRelease      string
	KubernetesVersion string
)

// DefaultOptions returns default settings.
func DefaultOptions() Options {
	return Options{
		BootstrapClusterName:    "sfyra",
		BootstrapTalosVmlinuz:   fmt.Sprintf("_out/%s/vmlinuz-amd64", TalosRelease),
		BootstrapTalosInitramfs: fmt.Sprintf("_out/%s/initramfs-amd64.xz", TalosRelease),
		BootstrapTalosInstaller: fmt.Sprintf("ghcr.io/siderolabs/installer:%s", TalosRelease),
		BootstrapCNIBundleURL:   fmt.Sprintf("https://github.com/siderolabs/talos/releases/download/%s/talosctl-cni-bundle-%s.tar.gz", TalosRelease, "amd64"),
		BootstrapCIDR:           "172.24.0.0/24",

		TalosKernelURL: fmt.Sprintf("https://github.com/siderolabs/talos/releases/download/%s/vmlinuz-amd64", TalosRelease),
		TalosInitrdURL: fmt.Sprintf("https://github.com/siderolabs/talos/releases/download/%s/initramfs-amd64.xz", TalosRelease),

		CoreProvider:            "cluster-api",
		BootstrapProviders:      []string{"talos"},
		InfrastructureProviders: []string{"sidero"},
		ControlPlaneProviders:   []string{"talos"},

		ManagementNodes: 4,

		BootstrapMemMB:  3072,
		BootstrapCPUs:   3,
		BootstrapDiskGB: 6,

		ManagementMemMB:  2048,
		ManagementCPUs:   2,
		ManagementDiskGB: 6,

		DefaultBootOrder: "cn", // disk, then network; override to "nc" to force PXE boot each time

		TalosctlPath: fmt.Sprintf("_out/%s/talosctl-linux-amd64", TalosRelease),
	}
}
