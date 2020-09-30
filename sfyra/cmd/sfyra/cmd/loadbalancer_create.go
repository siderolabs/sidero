// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"context"
	"net"

	"github.com/spf13/cobra"
	"github.com/talos-systems/talos/pkg/cli"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/talos-systems/sidero/sfyra/pkg/capi"
	"github.com/talos-systems/sidero/sfyra/pkg/loadbalancer"
	"github.com/talos-systems/sidero/sfyra/pkg/vm"
)

var (
	kubeconfig  string
	clusterName string
	lbPort      int
)

var loadbalancerCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a load balancer for the control plane nodes of Sidero cluster.",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cli.WithContext(context.Background(), func(ctx context.Context) error {
			bootSourceIP := net.ParseIP(bootSource)

			managementSet, err := vm.NewSet(ctx, vm.Options{
				Name:       options.ManagementSetName,
				Nodes:      options.ManagementNodes,
				BootSource: bootSourceIP,
				CIDR:       options.ManagementCIDR,

				TalosctlPath: options.TalosctlPath,

				CPUs:   options.CPUs,
				MemMB:  options.MemMB,
				DiskGB: options.DiskGB,
			})
			if err != nil {
				return err
			}

			if err = managementSet.Setup(ctx); err != nil {
				return err
			}

			config, err := clientcmd.BuildConfigFromKubeconfigGetter("", func() (*clientcmdapi.Config, error) {
				return clientcmd.LoadFromFile(kubeconfig)
			})
			if err != nil {
				return err
			}

			metalClient, err := capi.GetMetalClient(config)
			if err != nil {
				return err
			}

			lb, err := loadbalancer.NewControlPlane(metalClient, managementSet.BridgeIP(), lbPort, "default", clusterName, managementSet.Nodes())
			if err != nil {
				return err
			}

			defer lb.Close()

			<-ctx.Done()

			return nil
		})
	},
}

func init() {
	loadbalancerCmd.AddCommand(loadbalancerCreateCmd)

	loadbalancerCreateCmd.Flags().StringVar(&options.ManagementSetName, "management-set-name", options.ManagementSetName, "name for the management VM set")
	loadbalancerCreateCmd.Flags().IntVar(&options.ManagementNodes, "management-nodes", options.ManagementNodes, "number of PXE nodes to create for the management rack")
	loadbalancerCreateCmd.Flags().StringVar(&options.ManagementCIDR, "management-cidr", options.ManagementCIDR, "management cluster network CIDR")
	loadbalancerCreateCmd.Flags().StringVar(&bootSource, "boot-source", "172.24.0.2", "the boot source IP for the iPXE boot")
	loadbalancerCreateCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "path to kubeconfig for the CAPI cluster")
	loadbalancerCreateCmd.Flags().StringVar(&clusterName, "cluster-name", "management-cluster", "name of the cluster to build the loadbalancer for")
	loadbalancerCreateCmd.Flags().IntVar(&lbPort, "load-balancer-port", 16443, "port for the loadbalancer")
}
