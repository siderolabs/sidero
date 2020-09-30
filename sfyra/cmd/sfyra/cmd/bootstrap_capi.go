// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/talos-systems/talos/pkg/cli"

	"github.com/talos-systems/sidero/sfyra/pkg/bootstrap"
	"github.com/talos-systems/sidero/sfyra/pkg/capi"
)

var bootstrapCAPICmd = &cobra.Command{
	Use:   "capi",
	Short: "Install and patch CAPI.",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cli.WithContext(context.Background(), func(ctx context.Context) error {
			bootstrapCluster, err := bootstrap.NewCluster(ctx, bootstrap.Options{
				Name: options.BootstrapClusterName,
				CIDR: options.BootstrapCIDR,

				Vmlinuz:        options.BootstrapTalosVmlinuz,
				Initramfs:      options.BootstrapTalosInitramfs,
				InstallerImage: options.BootstrapTalosInstaller,

				TalosctlPath: options.TalosctlPath,

				RegistryMirrors: options.RegistryMirrors,

				CPUs:   options.CPUs,
				MemMB:  options.MemMB,
				DiskGB: options.DiskGB,
			})
			if err != nil {
				return err
			}

			if err = bootstrapCluster.Setup(ctx); err != nil {
				return err
			}

			clusterAPI, err := capi.NewManager(ctx, bootstrapCluster, capi.Options{
				ClusterctlConfigPath:    options.ClusterctlConfigPath,
				BootstrapProviders:      options.BootstrapProviders,
				InfrastructureProviders: options.InfrastructureProviders,
				ControlPlaneProviders:   options.ControlPlaneProviders,
			})
			if err != nil {
				return err
			}

			return clusterAPI.Install(ctx)
		})
	},
}

func init() {
	bootstrapCmd.AddCommand(bootstrapCAPICmd)

	bootstrapCAPICmd.Flags().StringVar(&options.BootstrapClusterName, "bootstrap-cluster-name", options.BootstrapClusterName, "bootstrap cluster name")
	bootstrapCAPICmd.Flags().StringVar(&options.BootstrapTalosVmlinuz, "bootstrap-vmlinuz", options.BootstrapTalosVmlinuz, "Talos kernel image for bootstrap cluster")
	bootstrapCAPICmd.Flags().StringVar(&options.BootstrapTalosInitramfs, "bootstrap-initramfs", options.BootstrapTalosInitramfs, "Talos initramfs image for bootstrap cluster")
	bootstrapCAPICmd.Flags().StringVar(&options.BootstrapTalosInstaller, "bootstrap-installer", options.BootstrapTalosInstaller, "Talos install image for bootstrap cluster")
	bootstrapCAPICmd.Flags().StringVar(&options.BootstrapCIDR, "bootstrap-cidr", options.BootstrapCIDR, "bootstrap cluster network CIDR")
	bootstrapCAPICmd.Flags().StringVar(&options.TalosctlPath, "talosctl-path", options.TalosctlPath, "path to the talosctl (for the QEMU provisioner)")
	bootstrapCAPICmd.Flags().StringSliceVar(&options.RegistryMirrors, "registry-mirrors", options.RegistryMirrors, "registry mirrors to use")
	bootstrapCAPICmd.Flags().StringVar(&options.ClusterctlConfigPath, "clusterctl-config", options.ClusterctlConfigPath, "path to the clusterctl config file")
}
