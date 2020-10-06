// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/talos-systems/talos/pkg/cli"

	"github.com/talos-systems/sidero/sfyra/pkg/bootstrap"
	"github.com/talos-systems/sidero/sfyra/pkg/capi"
	"github.com/talos-systems/sidero/sfyra/pkg/tests"
	"github.com/talos-systems/sidero/sfyra/pkg/vm"
)

var testIntegrationCmd = &cobra.Command{
	Use:   "integration",
	Short: "Run integration test against Sidero.",
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

				CPUs:   options.BootstrapCPUs,
				MemMB:  options.BootstrapMemMB,
				DiskGB: options.BootstrapDiskGB,
			})
			if err != nil {
				return err
			}

			if !options.SkipTeardown {
				defer bootstrapCluster.TearDown(ctx) //nolint: errcheck
			}

			if err = bootstrapCluster.Setup(ctx); err != nil {
				return err
			}

			managementSet, err := vm.NewSet(ctx, vm.Options{
				Name:       options.ManagementSetName,
				Nodes:      options.ManagementNodes,
				BootSource: bootstrapCluster.SideroComponentsIP(),
				CIDR:       options.ManagementCIDR,

				TalosctlPath: options.TalosctlPath,

				CPUs:   options.ManagementCPUs,
				MemMB:  options.ManagementMemMB,
				DiskGB: options.ManagementDiskGB,
			})
			if err != nil {
				return err
			}

			if !options.SkipTeardown {
				defer managementSet.TearDown(ctx) //nolint: errcheck
			}

			if err = managementSet.Setup(ctx); err != nil {
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

			if err = clusterAPI.Install(ctx); err != nil {
				return err
			}

			// hacky hack
			os.Args = append(os.Args[0:1], "-test.v")

			if ok := tests.Run(ctx, bootstrapCluster, managementSet, clusterAPI, tests.Options{
				KernelURL:      options.TalosKernelURL,
				InitrdURL:      options.TalosInitrdURL,
				InstallerImage: options.TalosInstaller,

				RegistryMirrors: options.RegistryMirrors,
			}); !ok {
				return fmt.Errorf("test failure")
			}

			return nil
		})
	},
}

func init() {
	testCmd.AddCommand(testIntegrationCmd)

	testIntegrationCmd.Flags().BoolVar(&options.SkipTeardown, "skip-teardown", options.SkipTeardown, "skip tearing down cluster")
	testIntegrationCmd.Flags().StringVar(&options.BootstrapClusterName, "bootstrap-cluster-name", options.BootstrapClusterName, "bootstrap cluster name")
	testIntegrationCmd.Flags().StringVar(&options.BootstrapTalosVmlinuz, "bootstrap-vmlinuz", options.BootstrapTalosVmlinuz, "Talos kernel image for bootstrap cluster")
	testIntegrationCmd.Flags().StringVar(&options.BootstrapTalosInitramfs, "bootstrap-initramfs", options.BootstrapTalosInitramfs, "Talos initramfs image for bootstrap cluster")
	testIntegrationCmd.Flags().StringVar(&options.BootstrapTalosInstaller, "bootstrap-installer", options.BootstrapTalosInstaller, "Talos install image for bootstrap cluster")
	testIntegrationCmd.Flags().StringVar(&options.BootstrapCIDR, "bootstrap-cidr", options.BootstrapCIDR, "bootstrap cluster network CIDR")
	testIntegrationCmd.Flags().StringVar(&options.ManagementCIDR, "management-cidr", options.ManagementCIDR, "management cluster network CIDR")
	testIntegrationCmd.Flags().IntVar(&options.ManagementNodes, "management-nodes", options.ManagementNodes, "number of PXE nodes to create for the management rack")
	testIntegrationCmd.Flags().StringVar(&options.TalosctlPath, "talosctl-path", options.TalosctlPath, "path to the talosctl (for the QEMU provisioner)")
	testIntegrationCmd.Flags().StringSliceVar(&options.RegistryMirrors, "registry-mirrors", options.RegistryMirrors, "registry mirrors to use")
	testIntegrationCmd.Flags().StringVar(&options.TalosKernelURL, "talos-kernel-url", options.TalosKernelURL, "Talos kernel image URL for Cluster API Environment")
	testIntegrationCmd.Flags().StringVar(&options.TalosInitrdURL, "talos-initrd-url", options.TalosInitrdURL, "Talos initramfs image URL for Cluster API Environment")
	testIntegrationCmd.Flags().StringVar(&options.ClusterctlConfigPath, "clusterctl-config", options.ClusterctlConfigPath, "path to the clusterctl config file")
}
