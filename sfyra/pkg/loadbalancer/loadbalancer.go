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
	"strconv"
	"sync"
	"time"

	cacpt "github.com/talos-systems/cluster-api-control-plane-provider-talos/api/v1alpha3"
	"github.com/talos-systems/go-loadbalancer/loadbalancer"
	"github.com/talos-systems/talos/pkg/provision"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"

	sidero "github.com/talos-systems/sidero/app/cluster-api-provider-sidero/api/v1alpha3"
	metal "github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
)

// ControlPlane implements dynamic loadbalancer for the control plane.
type ControlPlane struct {
	client client.Client

	endpoint string
	lb       loadbalancer.TCP

	prevUpstreams []string

	clusterNamespace, clusterName string
	nodes                         []provision.NodeInfo

	ctx       context.Context
	ctxCancel context.CancelFunc

	wg sync.WaitGroup
}

// NewControlPlane initializes new control plane load balancer.
func NewControlPlane(client client.Client, address net.IP, port int, clusterNamespace, clusterName string, nodes []provision.NodeInfo, verboseLog bool) (*ControlPlane, error) {
	cp := ControlPlane{
		client:           client,
		clusterNamespace: clusterNamespace,
		clusterName:      clusterName,
		nodes:            nodes,
	}

	cp.ctx, cp.ctxCancel = context.WithCancel(context.Background())

	var err error

	if port == 0 {
		port, err = findListenPort(address)
		if err != nil {
			return nil, err
		}
	}

	cp.endpoint = net.JoinHostPort(address.String(), strconv.Itoa(port))

	if !verboseLog {
		// send logs to /dev/null
		cp.lb.Logger = log.New(ioutil.Discard, "", 0)
	}

	// create route without any upstreams yet
	if err := cp.lb.AddRoute(cp.endpoint, nil); err != nil {
		return nil, err
	}

	cp.wg.Add(1)

	go cp.reconcileLoop()

	return &cp, cp.lb.Start()
}

// GetEndpoint returns loadbalancer endpoint.
func (cp *ControlPlane) GetEndpoint() string {
	return cp.endpoint
}

// Close the load balancer.
func (cp *ControlPlane) Close() error {
	cp.ctxCancel()
	cp.wg.Wait()

	if err := cp.lb.Close(); err != nil {
		return err
	}

	return cp.lb.Wait()
}

func (cp *ControlPlane) reconcileLoop() {
	defer cp.wg.Done()

	const interval = 15 * time.Second

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		if err := cp.reconcile(); err != nil {
			log.Printf("load balancer reconcile failed: %s", err)
		}

		select {
		case <-cp.ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (cp *ControlPlane) reconcile() error {
	var cluster v1alpha3.Cluster

	if err := cp.client.Get(cp.ctx, types.NamespacedName{Namespace: cp.clusterNamespace, Name: cp.clusterName}, &cluster); err != nil {
		return err
	}

	var controlPlane cacpt.TalosControlPlane

	if err := cp.client.Get(cp.ctx, types.NamespacedName{Namespace: cluster.Spec.ControlPlaneRef.Namespace, Name: cluster.Spec.ControlPlaneRef.Name}, &controlPlane); err != nil {
		return err
	}

	var machines v1alpha3.MachineList

	labelSelector, err := labels.Parse(controlPlane.Status.Selector)
	if err != nil {
		return err
	}

	if err := cp.client.List(cp.ctx, &machines, client.MatchingLabelsSelector{Selector: labelSelector}); err != nil {
		return err
	}

	var upstreams []string

	for _, machine := range machines.Items {
		var metalMachine sidero.MetalMachine

		if err := cp.client.Get(cp.ctx, types.NamespacedName{Namespace: machine.Spec.InfrastructureRef.Namespace, Name: machine.Spec.InfrastructureRef.Name}, &metalMachine); err != nil {
			continue
		}

		var server metal.Server

		if metalMachine.Spec.ServerRef == nil {
			continue
		}

		if err := cp.client.Get(cp.ctx, types.NamespacedName{Namespace: metalMachine.Spec.ServerRef.Namespace, Name: metalMachine.Spec.ServerRef.Name}, &server); err != nil {
			return err
		}

		for _, node := range cp.nodes {
			if node.UUID.String() == server.Name {
				upstreams = append(upstreams, fmt.Sprintf("%s:6443", node.PrivateIP.String()))

				break
			}
		}
	}

	sort.Strings(upstreams)

	if !reflect.DeepEqual(cp.prevUpstreams, upstreams) {
		log.Printf("new control plane loadbalancer %q routes: %v", cp.endpoint, upstreams)
	}

	cp.prevUpstreams = upstreams

	return cp.lb.ReconcileRoute(cp.endpoint, upstreams)
}

func findListenPort(address net.IP) (int, error) {
	l, err := net.Listen("tcp", net.JoinHostPort(address.String(), "0"))
	if err != nil {
		return 0, err
	}

	port := l.Addr().(*net.TCPAddr).Port

	return port, l.Close()
}
