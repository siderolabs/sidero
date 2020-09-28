// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"testing"

	"github.com/talos-systems/talos/pkg/cli"

	"github.com/talos-systems/sidero/sfyra/pkg/bootstrap"
	"github.com/talos-systems/sidero/sfyra/pkg/capi"
	"github.com/talos-systems/sidero/sfyra/pkg/tests"
	"github.com/talos-systems/sidero/sfyra/pkg/vm"
)

func main() {
	options := DefaultOptions()

	flag.BoolVar(&options.SkipTeardown, "skip-teardown", options.SkipTeardown, "skip tearing down cluster")
	flag.StringVar(&options.BootstrapClusterName, "bootstrap-cluster-name", options.BootstrapClusterName, "bootstrap cluster name")
	flag.StringVar(&options.BootstrapTalosVmlinuz, "bootstrap-vmlinuz", options.BootstrapTalosVmlinuz, "Talos kernel image for bootstrap cluster")
	flag.StringVar(&options.BootstrapTalosInitramfs, "bootstrap-initramfs", options.BootstrapTalosInitramfs, "Talos initramfs image for bootstrap cluster")
	flag.StringVar(&options.BootstrapTalosInstaller, "bootstrap-installer", options.BootstrapTalosInstaller, "Talos install image for bootstrap cluster")
	flag.StringVar(&options.BootstrapCIDR, "bootstrap-cidr", options.BootstrapCIDR, "bootstrap cluster network CIDR")
	flag.StringVar(&options.ManagementCIDR, "management-cidr", options.ManagementCIDR, "management cluster network CIDR")
	flag.IntVar(&options.ManagementNodes, "management-nodes", options.ManagementNodes, "number of PXE nodes to create for the management rack")
	flag.StringVar(&options.TalosctlPath, "talosctl-path", options.TalosctlPath, "path to the talosctl (for qemu provisioner)")
	flag.Var(&options.RegistryMirrors, "registry-mirrors", "registry mirrors to use")
	flag.StringVar(&options.TalosKernelURL, "talos-kernel-url", options.TalosKernelURL, "Talos kernel image URL for Cluster API Environment")
	flag.StringVar(&options.TalosInitrdURL, "talos-initrd-url", options.TalosInitrdURL, "Talos initramfs image URL for Cluster API Environment")
	flag.StringVar(&options.ClusterctlConfigPath, "clusterctl-config", options.ClusterctlConfigPath, "path to the clusterctl config file")

	testing.Init()

	flag.Parse()

	err := cli.WithContext(context.Background(), func(ctx context.Context) error {
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

		if !options.SkipTeardown {
			defer bootstrapCluster.TearDown(ctx) //nolint: errcheck
		}

		if err = bootstrapCluster.Setup(ctx); err != nil {
			return err
		}

		managementSet, err := vm.NewSet(ctx, vm.Options{
			Name:       options.BootstrapClusterName + "-management",
			Nodes:      options.ManagementNodes,
			BootSource: bootstrapCluster.SideroComponentsIP(),
			CIDR:       options.ManagementCIDR,

			TalosctlPath: options.TalosctlPath,

			CPUs:   options.CPUs,
			MemMB:  options.MemMB,
			DiskGB: options.DiskGB,
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
	if err != nil {
		log.Fatal(err)
	}
}
