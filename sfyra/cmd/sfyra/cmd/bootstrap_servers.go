// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"context"
	"net"

	"github.com/spf13/cobra"
	"github.com/talos-systems/talos/pkg/cli"

	"github.com/talos-systems/sidero/sfyra/pkg/vm"
)

var bootSource string

var bootstrapServersCmd = &cobra.Command{
	Use:   "servers",
	Short: "Create a set of VMs ready for PXE booting.",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cli.WithContext(context.Background(), func(ctx context.Context) error {
			bootSourceIP := net.ParseIP(bootSource)

			vmSet, err := vm.NewSet(ctx, vm.Options{
				Name:       options.ManagementSetName,
				Nodes:      options.ManagementNodes,
				BootSource: bootSourceIP,
				CIDR:       options.ManagementCIDR,

				TalosctlPath: options.TalosctlPath,

				CPUs:   options.ManagementCPUs,
				MemMB:  options.ManagementMemMB,
				DiskGB: options.ManagementDiskGB,

				DefaultBootOrder: options.DefaultBootOrder,
			})
			if err != nil {
				return err
			}

			return vmSet.Setup(ctx)
		})
	},
}

func init() {
	bootstrapCmd.AddCommand(bootstrapServersCmd)

	bootstrapServersCmd.Flags().StringVar(&options.ManagementSetName, "management-set-name", options.ManagementSetName, "name for the management VM set")
	bootstrapServersCmd.Flags().IntVar(&options.ManagementNodes, "management-nodes", options.ManagementNodes, "number of PXE nodes to create for the management rack")
	bootstrapServersCmd.Flags().StringVar(&options.ManagementCIDR, "management-cidr", options.ManagementCIDR, "management cluster network CIDR")
	bootstrapServersCmd.Flags().StringVar(&bootSource, "boot-source", "172.24.0.2", "the boot source IP for the iPXE boot")
	bootstrapServersCmd.Flags().StringVar(&options.DefaultBootOrder, "default-boot-order", options.DefaultBootOrder, "QEMU default boot order")
}
