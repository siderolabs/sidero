// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package loadbalancer provides dynamic loadbalancer for the controlplane of the CAPI cluster.
package loadbalancer

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"reflect"
	"sort"
	"sync"
	"time"

	cacpt "github.com/talos-systems/cluster-api-control-plane-provider-talos/api/v1alpha3"
	"github.com/talos-systems/go-loadbalancer/controlplane"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	capiv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "github.com/talos-systems/sidero/app/caps-controller-manager/api/v1alpha3"
	metal "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha1"
)

// ControlPlane implements dynamic loadbalancer for the control plane.
type ControlPlane struct {
	client client.Client

	prevUpstreams []string

	clusterNamespace, clusterName string

	lb *controlplane.LoadBalancer

	ctxCancel context.CancelFunc

	wg sync.WaitGroup
}

// NewControlPlane initializes new control plane load balancer.
func NewControlPlane(client client.Client, address net.IP, port int, clusterNamespace, clusterName string, verboseLog bool) (*ControlPlane, error) {
	cp := ControlPlane{
		client:           client,
		clusterNamespace: clusterNamespace,
		clusterName:      clusterName,
	}

	logWriter := log.Writer()
	if !verboseLog {
		logWriter = ioutil.Discard
	}

	var err error

	cp.lb, err = controlplane.NewLoadBalancer(address.String(), port, logWriter)
	if err != nil {
		return nil, err
	}

	upstreamCh := make(chan []string)

	var ctx context.Context

	ctx, cp.ctxCancel = context.WithCancel(context.Background())

	cp.wg.Add(1)

	go cp.reconcileLoop(ctx, upstreamCh)

	return &cp, cp.lb.Start(upstreamCh)
}

// GetEndpoint returns loadbalancer endpoint.
func (cp *ControlPlane) GetEndpoint() string {
	return cp.lb.Endpoint()
}

// Close the load balancer.
func (cp *ControlPlane) Close() error {
	cp.ctxCancel()
	cp.wg.Wait()

	return cp.lb.Shutdown()
}

func (cp *ControlPlane) reconcileLoop(ctx context.Context, upstreamCh chan<- []string) {
	defer cp.wg.Done()

	const interval = 15 * time.Second

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		if err := cp.reconcile(ctx); err != nil {
			log.Printf("load balancer reconcile failed: %s", err)
		} else {
			select {
			case upstreamCh <- cp.prevUpstreams:
			case <-ctx.Done():
				return
			}
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (cp *ControlPlane) reconcile(ctx context.Context) error {
	var cluster capiv1.Cluster

	if err := cp.client.Get(ctx, types.NamespacedName{Namespace: cp.clusterNamespace, Name: cp.clusterName}, &cluster); err != nil {
		return err
	}

	var controlPlane cacpt.TalosControlPlane

	if err := cp.client.Get(ctx, types.NamespacedName{Namespace: cluster.Spec.ControlPlaneRef.Namespace, Name: cluster.Spec.ControlPlaneRef.Name}, &controlPlane); err != nil {
		return err
	}

	var machines capiv1.MachineList

	labelSelector, err := labels.Parse(controlPlane.Status.Selector)
	if err != nil {
		return err
	}

	if err := cp.client.List(ctx, &machines, client.MatchingLabelsSelector{Selector: labelSelector}); err != nil {
		return err
	}

	var upstreams []string

	for _, machine := range machines.Items {
		// we could have looked up addresses via Machine.status, but as we still have tests with Talos 0.13 (before SideroLink was introduced),
		// we need to keep this way of looking up addresses.
		var metalMachine infrav1.MetalMachine

		if err := cp.client.Get(ctx, types.NamespacedName{Namespace: machine.Spec.InfrastructureRef.Namespace, Name: machine.Spec.InfrastructureRef.Name}, &metalMachine); err != nil {
			continue
		}

		var server metal.Server

		if metalMachine.Spec.ServerRef == nil {
			continue
		}

		if err := cp.client.Get(ctx, types.NamespacedName{Namespace: metalMachine.Spec.ServerRef.Namespace, Name: metalMachine.Spec.ServerRef.Name}, &server); err != nil {
			return err
		}

		for _, address := range server.Status.Addresses {
			if address.Type == corev1.NodeInternalIP {
				upstreams = append(upstreams, fmt.Sprintf("%s:6443", address.Address))
			}
		}
	}

	sort.Strings(upstreams)

	if !reflect.DeepEqual(cp.prevUpstreams, upstreams) {
		log.Printf("new control plane loadbalancer %q routes: %v", cp.lb.Endpoint(), upstreams)
	}

	cp.prevUpstreams = upstreams

	return nil
}
