// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package siderolink

import (
	"context"
	"net/netip"
	"time"

	"go.uber.org/zap"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	sidero "github.com/siderolabs/sidero/app/caps-controller-manager/api/v1alpha3"
	"github.com/siderolabs/siderolink/pkg/wireguard"
)

// PeerState syncs data from Kubernetes ServerBinding as peer state.
type PeerState struct {
	kubeconfig *rest.Config
	logger     *zap.Logger

	eventCh chan wireguard.PeerEvent
}

// NewPeerState initializes PeerState.
func NewPeerState(kubeconfig *rest.Config, logger *zap.Logger) *PeerState {
	return &PeerState{
		kubeconfig: kubeconfig,
		logger:     logger,
		eventCh:    make(chan wireguard.PeerEvent, 16),
	}
}

// Run the watch loop reporting peer state changes.
//
//nolint:dupl
func (peers *PeerState) Run(ctx context.Context) error {
	dc, err := dynamic.NewForConfig(peers.kubeconfig)
	if err != nil {
		return err
	}

	// Create a factory object that can generate informers for resource types
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dc, 10*time.Minute, "", nil)

	informerFactory := factory.ForResource(sidero.GroupVersion.WithResource("serverbindings"))
	informer := informerFactory.Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(new interface{}) { peers.notify(nil, new) },
		UpdateFunc: peers.notify,
		DeleteFunc: func(old interface{}) { peers.notify(old, nil) },
	})

	informer.Run(ctx.Done())

	return nil
}

func (peers *PeerState) notify(old, new interface{}) {
	var oldServerBinding, newServerBinding *sidero.ServerBinding

	if old != nil {
		oldServerBinding = &sidero.ServerBinding{}

		err := runtime.DefaultUnstructuredConverter.
			FromUnstructured(old.(*unstructured.Unstructured).UnstructuredContent(), oldServerBinding)
		if err != nil {
			peers.logger.Error("failed converting old event object", zap.Error(err))

			return
		}
	}

	if new != nil {
		newServerBinding = &sidero.ServerBinding{}

		err := runtime.DefaultUnstructuredConverter.
			FromUnstructured(new.(*unstructured.Unstructured).UnstructuredContent(), newServerBinding)
		if err != nil {
			peers.logger.Error("failed converting new event object", zap.Error(err))

			return
		}
	}

	if oldServerBinding != nil && newServerBinding != nil {
		if oldServerBinding.Spec.SideroLink == newServerBinding.Spec.SideroLink {
			// no change to SideroLink, skip it
			return
		}
	}

	if oldServerBinding != nil {
		peers.buildEvent(oldServerBinding, true)
	}

	if newServerBinding != nil {
		peers.buildEvent(newServerBinding, false)
	}
}

func (peers *PeerState) buildEvent(serverBinding *sidero.ServerBinding, deleted bool) {
	if serverBinding.Spec.SideroLink.NodePublicKey == "" || serverBinding.Spec.SideroLink.NodeAddress == "" {
		// no SideroLink information
		return
	}

	address, err := netip.ParsePrefix(serverBinding.Spec.SideroLink.NodeAddress)
	if err != nil {
		peers.logger.Error("error parsing node address", zap.Error(err), zap.String("uuid", serverBinding.Name))

		return
	}

	pubKey, err := wgtypes.ParseKey(serverBinding.Spec.SideroLink.NodePublicKey)
	if err != nil {
		peers.logger.Error("error parsing public key", zap.Error(err), zap.String("uuid", serverBinding.Name))

		return
	}

	peers.eventCh <- wireguard.PeerEvent{
		PubKey:  pubKey,
		Remove:  deleted,
		Address: address.Addr(),
	}
}

// EventCh implements the wireguard.PeerSource interface.
func (peers *PeerState) EventCh() <-chan wireguard.PeerEvent {
	return peers.eventCh
}
