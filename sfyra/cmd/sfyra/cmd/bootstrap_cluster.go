// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"context"

	"github.com/siderolabs/talos/pkg/cli"
	"github.com/spf13/cobra"

	"github.com/siderolabs/sidero/sfyra/pkg/bootstrap"
)

var bootstrapClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Create a Talos cluster to be used as bootstrap Sidero cluster.",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cli.WithContext(context.Background(), func(ctx context.Context) error {
			bootstrapCluster, err := bootstrap.NewCluster(ctx, bootstrap.Options{
				Name: options.BootstrapClusterName,
				CIDR: options.BootstrapCIDR,

				Vmlinuz:        options.BootstrapTalosVmlinuz,
				Initramfs:      options.BootstrapTalosInitramfs,
				InstallerImage: options.BootstrapTalosInstaller,
				CNIBundleURL:   options.BootstrapCNIBundleURL,

				TalosctlPath: options.TalosctlPath,

				RegistryMirrors: options.RegistryMirrors,

				BootstrapCPUs:   options.BootstrapCPUs,
				BootstrapMemMB:  options.BootstrapMemMB,
				BootstrapDiskGB: options.BootstrapDiskGB,
			})
			if err != nil {
				return err
			}

			return bootstrapCluster.Setup(ctx)
		})
	},
}

func init() {
	bootstrapCmd.AddCommand(bootstrapClusterCmd)

	bootstrapClusterCmd.Flags().StringVar(&options.BootstrapClusterName, "bootstrap-cluster-name", options.BootstrapClusterName, "bootstrap cluster name")
	bootstrapClusterCmd.Flags().StringVar(&options.BootstrapTalosVmlinuz, "bootstrap-vmlinuz", options.BootstrapTalosVmlinuz, "Talos kernel image for bootstrap cluster")
	bootstrapClusterCmd.Flags().StringVar(&options.BootstrapTalosInitramfs, "bootstrap-initramfs", options.BootstrapTalosInitramfs, "Talos initramfs image for bootstrap cluster")
	bootstrapClusterCmd.Flags().StringVar(&options.BootstrapTalosInstaller, "bootstrap-installer", options.BootstrapTalosInstaller, "Talos install image for bootstrap cluster")
	bootstrapClusterCmd.Flags().StringVar(&options.BootstrapCIDR, "bootstrap-cidr", options.BootstrapCIDR, "bootstrap cluster network CIDR")
	bootstrapClusterCmd.Flags().StringVar(&options.TalosctlPath, "talosctl-path", options.TalosctlPath, "path to the talosctl (for the QEMU provisioner)")
	bootstrapClusterCmd.Flags().StringSliceVar(&options.RegistryMirrors, "registry-mirror", options.RegistryMirrors, "registry mirrors to use")
	bootstrapClusterCmd.Flags().StringSliceVar(&options.RegistryMirrors, "registry-mirrors", options.RegistryMirrors, "registry mirrors to use")
	Should(bootstrapClusterCmd.Flags().MarkDeprecated("registry-mirrors", "please use --registry-mirror (singular) instead"))
}
