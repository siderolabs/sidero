// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"context"
	"log"
	"net/netip"

	talosnet "github.com/siderolabs/net"
	"github.com/siderolabs/talos/pkg/cli"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/siderolabs/sidero/sfyra/pkg/capi"
	"github.com/siderolabs/sidero/sfyra/pkg/loadbalancer"
)

var (
	kubeconfig  string
	clusterName string
	lbAddress   string
	lbPort      int
)

var loadbalancerCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a load balancer for the control plane nodes of Sidero cluster.",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cli.WithContext(context.Background(), func(ctx context.Context) error {
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

			lb, err := loadbalancer.NewControlPlane(metalClient, netip.MustParseAddr(lbAddress), lbPort, "default", clusterName, true)
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

	cidr := netip.MustParsePrefix(options.BootstrapCIDR)

	bridgeIP, err := talosnet.NthIPInNetwork(cidr, 1)
	if err != nil {
		log.Fatal(err)
	}

	loadbalancerCreateCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "path to kubeconfig for the CAPI cluster")
	loadbalancerCreateCmd.Flags().StringVar(&clusterName, "cluster-name", "management-cluster", "name of the cluster to build the loadbalancer for")
	loadbalancerCreateCmd.Flags().StringVar(&lbAddress, "load-balancer-address", bridgeIP.String(), "address for the loadbalancer")
	loadbalancerCreateCmd.Flags().IntVar(&lbPort, "load-balancer-port", 16443, "port for the loadbalancer")
}
