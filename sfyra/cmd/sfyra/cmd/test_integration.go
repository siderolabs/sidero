// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/siderolabs/talos/pkg/cli"
	"github.com/spf13/cobra"

	"github.com/siderolabs/sidero/sfyra/pkg/bootstrap"
	"github.com/siderolabs/sidero/sfyra/pkg/capi"
	"github.com/siderolabs/sidero/sfyra/pkg/tests"
)

var runTestPattern string

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
				CNIBundleURL:   options.BootstrapCNIBundleURL,

				TalosctlPath: options.TalosctlPath,

				RegistryMirrors: options.RegistryMirrors,

				BootstrapCPUs:   options.BootstrapCPUs,
				BootstrapMemMB:  options.BootstrapMemMB,
				BootstrapDiskGB: options.BootstrapDiskGB,

				VMNodes: options.ManagementNodes,

				VMCPUs:   options.ManagementCPUs,
				VMMemMB:  options.ManagementMemMB,
				VMDiskGB: options.ManagementDiskGB,

				VMDefaultBootOrder: options.DefaultBootOrder,
			})
			if err != nil {
				return err
			}

			if !options.SkipTeardown {
				defer bootstrapCluster.TearDown(ctx) //nolint:errcheck
			}

			if err = bootstrapCluster.Setup(ctx); err != nil {
				return err
			}

			clusterAPI, err := capi.NewManager(ctx, bootstrapCluster, capi.Options{
				ClusterctlConfigPath:    options.ClusterctlConfigPath,
				CoreProvider:            options.CoreProvider,
				BootstrapProviders:      options.BootstrapProviders,
				InfrastructureProviders: options.InfrastructureProviders,
				ControlPlaneProviders:   options.ControlPlaneProviders,

				PowerSimulatedExplicitFailureProb: options.PowerSimulatedExplicitFailureProb,
				PowerSimulatedSilentFailureProb:   options.PowerSimulatedSilentFailureProb,
			})
			if err != nil {
				return err
			}

			if err = clusterAPI.Install(ctx); err != nil {
				return err
			}

			// hacky hack
			os.Args = append(os.Args[0:1], "-test.v")

			if ok := tests.Run(ctx, bootstrapCluster, clusterAPI, tests.Options{
				KernelURL: options.TalosKernelURL,
				InitrdURL: options.TalosInitrdURL,

				RegistryMirrors: options.RegistryMirrors,

				RunTestPattern: runTestPattern,

				TalosRelease:      TalosRelease,
				KubernetesVersion: KubernetesVersion,
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
	testIntegrationCmd.Flags().IntVar(&options.ManagementNodes, "management-nodes", options.ManagementNodes, "number of PXE nodes to create for the management rack")
	testIntegrationCmd.Flags().StringVar(&options.TalosctlPath, "talosctl-path", options.TalosctlPath, "path to the talosctl (for the QEMU provisioner)")
	testIntegrationCmd.Flags().StringSliceVar(&options.RegistryMirrors, "registry-mirror", options.RegistryMirrors, "registry mirrors to use")
	testIntegrationCmd.Flags().StringSliceVar(&options.RegistryMirrors, "registry-mirrors", options.RegistryMirrors, "registry mirrors to use")
	Should(testIntegrationCmd.Flags().MarkDeprecated("registry-mirrors", "please use --registry-mirror (singular) instead"))
	testIntegrationCmd.Flags().StringVar(&options.TalosKernelURL, "talos-kernel-url", options.TalosKernelURL, "Talos kernel image URL for Cluster API Environment")
	testIntegrationCmd.Flags().StringVar(&options.TalosInitrdURL, "talos-initrd-url", options.TalosInitrdURL, "Talos initramfs image URL for Cluster API Environment")
	testIntegrationCmd.Flags().StringVar(&options.ClusterctlConfigPath, "clusterctl-config", options.ClusterctlConfigPath, "path to the clusterctl config file")
	testIntegrationCmd.Flags().StringVar(&options.DefaultBootOrder, "default-boot-order", options.DefaultBootOrder, "QEMU default boot order")
	testIntegrationCmd.Flags().Float64Var(&options.PowerSimulatedExplicitFailureProb, "power-simulated-explicit-failure-prob", options.PowerSimulatedExplicitFailureProb, "simulated power management explicit failure probability")
	testIntegrationCmd.Flags().Float64Var(&options.PowerSimulatedSilentFailureProb, "power-simulated-silent-failure-prob", options.PowerSimulatedSilentFailureProb, "simulated power management silent failure probability")
	testIntegrationCmd.Flags().StringVar(&runTestPattern, "test.run", "", "tests to run (regular expression)")
}
