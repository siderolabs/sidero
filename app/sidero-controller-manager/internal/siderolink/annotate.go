// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package siderolink

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"inet.af/netaddr"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	sidero "github.com/talos-systems/sidero/app/caps-controller-manager/api/v1alpha3"
)

// Annotation describes the source server by SideroLink IP address.
type Annotation struct {
	ServerUUID       string
	Namespace        string
	MetalMachineName string
	MachineName      string
	ClusterName      string
}

// Annotator keeps a cache of annotations per SideroLink IP address.
type Annotator struct {
	logger      *zap.Logger
	metalClient runtimeclient.Client
	kubeconfig  *rest.Config

	nodesMu sync.Mutex
	nodes   map[string]Annotation
}

// NewAnnotator initializes new server.
func NewAnnotator(metalClient runtimeclient.Client, kubeconfig *rest.Config, logger *zap.Logger) *Annotator {
	return &Annotator{
		logger:      logger,
		kubeconfig:  kubeconfig,
		metalClient: metalClient,
		nodes:       map[string]Annotation{},
	}
}

// Run the watch loop on ServerBindings to build the annotation database.
//
//nolint:dupl
func (a *Annotator) Run(ctx context.Context) error {
	dc, err := dynamic.NewForConfig(a.kubeconfig)
	if err != nil {
		return err
	}

	// Create a factory object that can generate informers for resource types
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dc, 10*time.Minute, "", nil)

	informerFactory := factory.ForResource(sidero.GroupVersion.WithResource("serverbindings"))
	informer := informerFactory.Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(new interface{}) { a.notify(nil, new) },
		UpdateFunc: a.notify,
		DeleteFunc: func(old interface{}) { a.notify(old, nil) },
	})

	informer.Run(ctx.Done())

	return nil
}

func (a *Annotator) Get(addr string) (Annotation, bool) {
	a.nodesMu.Lock()
	defer a.nodesMu.Unlock()

	annotation, exists := a.nodes[addr]

	return annotation, exists
}

func (a *Annotator) notify(old, new interface{}) {
	var oldServerBinding, newServerBinding *sidero.ServerBinding

	if old != nil {
		oldServerBinding = &sidero.ServerBinding{}

		err := runtime.DefaultUnstructuredConverter.
			FromUnstructured(old.(*unstructured.Unstructured).UnstructuredContent(), oldServerBinding)
		if err != nil {
			a.logger.Error("failed converting old event object", zap.Error(err))

			return
		}
	}

	if new != nil {
		newServerBinding = &sidero.ServerBinding{}

		err := runtime.DefaultUnstructuredConverter.
			FromUnstructured(new.(*unstructured.Unstructured).UnstructuredContent(), newServerBinding)
		if err != nil {
			a.logger.Error("failed converting new event object", zap.Error(err))

			return
		}
	}

	a.nodesMu.Lock()
	defer a.nodesMu.Unlock()

	if new == nil {
		delete(a.nodes, oldServerBinding.Spec.SideroLink.NodeAddress)
	} else {
		if oldServerBinding != nil && oldServerBinding.Spec.SideroLink.NodeAddress == newServerBinding.Spec.SideroLink.NodeAddress {
			// no change to the node address
			return
		}

		address := newServerBinding.Spec.SideroLink.NodeAddress
		if address == "" {
			return
		}

		ipPrefix, err := netaddr.ParseIPPrefix(address)
		if err != nil {
			a.logger.Error("failure parsing siderolink address", zap.Error(err))
			return
		}

		address = ipPrefix.IP().String()

		if oldServerBinding != nil {
			delete(a.nodes, oldServerBinding.Spec.SideroLink.NodeAddress)
		}

		annotation, err := a.buildAnnotation(newServerBinding)
		if err != nil {
			a.logger.Error("failure building annotation", zap.Error(err))
		}

		a.nodes[address] = annotation

		a.logger.Debug("new node mapping", zap.String("ip", address), zap.Any("annotation", annotation))
	}
}

func (a *Annotator) buildAnnotation(serverBinding *sidero.ServerBinding) (annotation Annotation, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	annotation.ServerUUID = serverBinding.Name
	annotation.Namespace = serverBinding.Spec.MetalMachineRef.Namespace
	annotation.MetalMachineName = serverBinding.Spec.MetalMachineRef.Name
	annotation.ClusterName = serverBinding.Labels[clusterv1.ClusterLabelName]

	var metalMachine sidero.MetalMachine

	if err = a.metalClient.Get(ctx,
		types.NamespacedName{
			Namespace: serverBinding.Spec.MetalMachineRef.Namespace,
			Name:      serverBinding.Spec.MetalMachineRef.Name,
		},
		&metalMachine); err != nil {
		return annotation, fmt.Errorf("error getting metal machine: %w", err)
	}

	for _, ref := range metalMachine.OwnerReferences {
		gv, err := schema.ParseGroupVersion(ref.APIVersion)
		if err != nil {
			continue
		}

		if ref.Kind == "Machine" && gv.Group == clusterv1.GroupVersion.Group {
			annotation.MachineName = ref.Name

			break
		}
	}

	return annotation, nil
}
