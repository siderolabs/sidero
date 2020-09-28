// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package capi

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	cabpt "github.com/talos-systems/cluster-api-bootstrap-provider-talos/api/v1alpha3"
	cacpt "github.com/talos-systems/cluster-api-control-plane-provider-talos/api/v1alpha3"
	taloscluster "github.com/talos-systems/talos/pkg/cluster"
	talosclusterapi "github.com/talos-systems/talos/pkg/machinery/api/cluster"
	talosclient "github.com/talos-systems/talos/pkg/machinery/client"
	clientconfig "github.com/talos-systems/talos/pkg/machinery/client/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/cluster-api/api/v1alpha3"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	sidero "github.com/talos-systems/sidero/app/cluster-api-provider-sidero/api/v1alpha3"
	"github.com/talos-systems/sidero/sfyra/pkg/vm"
)

// Cluster attaches to the provisioned CAPI cluster and provides talos.Cluster.
type Cluster struct {
	name        string
	endpoints   []string
	bridgeIP    net.IP
	client      *talosclient.Client
	k8sProvider *taloscluster.KubernetesClient
}

// NewCluster fetches cluster info from the CAPI state.
func NewCluster(ctx context.Context, metalClient runtimeclient.Reader, clusterName string, vmSet *vm.Set) (*Cluster, error) {
	var (
		cluster      v1alpha3.Cluster
		controlPlane cacpt.TalosControlPlane
		machines     v1alpha3.MachineList
		talosConfig  cabpt.TalosConfig
	)

	if err := metalClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: clusterName}, &cluster); err != nil {
		return nil, err
	}

	if err := metalClient.Get(ctx, types.NamespacedName{Namespace: cluster.Spec.ControlPlaneRef.Namespace, Name: cluster.Spec.ControlPlaneRef.Name}, &controlPlane); err != nil {
		return nil, err
	}

	labelSelector, err := labels.Parse(controlPlane.Status.Selector)
	if err != nil {
		return nil, err
	}

	if err = metalClient.List(ctx, &machines, runtimeclient.MatchingLabelsSelector{Selector: labelSelector}); err != nil {
		return nil, err
	}

	if len(machines.Items) < 1 {
		return nil, fmt.Errorf("not enough machines found")
	}

	configRef := machines.Items[0].Spec.Bootstrap.ConfigRef

	if err = metalClient.Get(ctx, types.NamespacedName{Namespace: configRef.Namespace, Name: configRef.Name}, &talosConfig); err != nil {
		return nil, err
	}

	var clientConfig *clientconfig.Config
	clientConfig, err = clientconfig.FromString(talosConfig.Status.TalosConfig)

	if err != nil {
		return nil, err
	}

	// TODO: endpoints in talosconfig should be filled by Sidero
	var metalMachine sidero.MetalMachine

	if err = metalClient.Get(ctx,
		types.NamespacedName{Namespace: machines.Items[0].Spec.InfrastructureRef.Namespace, Name: machines.Items[0].Spec.InfrastructureRef.Name},
		&metalMachine); err != nil {
		return nil, err
	}

	nodeUUID := metalMachine.Spec.ServerRef.Name

	for _, node := range vmSet.Nodes() {
		if node.UUID.String() == nodeUUID {
			clientConfig.Contexts[clientConfig.Context].Endpoints = append(clientConfig.Contexts[clientConfig.Context].Endpoints, node.PrivateIP.String())
		}
	}

	var talosClient *talosclient.Client

	talosClient, err = talosclient.New(ctx, talosclient.WithConfig(clientConfig))
	if err != nil {
		return nil, err
	}

	return &Cluster{
		name:      clusterName,
		endpoints: clientConfig.Contexts[clientConfig.Context].Endpoints,
		bridgeIP:  vmSet.BridgeIP(),
		client:    talosClient,
		k8sProvider: &taloscluster.KubernetesClient{
			ClientProvider: &taloscluster.ConfigClientProvider{
				DefaultClient: talosClient,
			},
		},
	}, nil
}

// Health runs the healthcheck for the cluster.
func (cluster *Cluster) Health(ctx context.Context) error {
	resp, err := cluster.client.ClusterHealthCheck(talosclient.WithNodes(ctx, cluster.endpoints...), 3*time.Minute, &talosclusterapi.ClusterInfo{})
	if err != nil {
		return err
	}

	if err := resp.CloseSend(); err != nil {
		return err
	}

	for {
		msg, err := resp.Recv()
		if err != nil {
			if err == io.EOF || status.Code(err) == codes.Canceled {
				return nil
			}

			return err
		}

		if msg.GetMetadata().GetError() != "" {
			return fmt.Errorf("healthcheck error: %s", msg.GetMetadata().GetError())
		}

		fmt.Fprintln(os.Stderr, msg.GetMessage())
	}
}

// Name of the cluster.
func (cluster *Cluster) Name() string {
	return cluster.name
}

// BridgeIP returns IP of the bridge which controls the cluster.
func (cluster *Cluster) BridgeIP() net.IP {
	return cluster.bridgeIP
}

// SideroComponentsIP returns the IP for the Sidero components (TFTP, iPXE, etc.).
func (cluster *Cluster) SideroComponentsIP() net.IP {
	panic("not implemented yet")
}

// KubernetesClient provides K8s client source.
func (cluster *Cluster) KubernetesClient() taloscluster.K8sProvider {
	return cluster.k8sProvider
}
