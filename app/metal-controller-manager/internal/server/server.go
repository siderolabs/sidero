// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"reflect"

	"google.golang.org/grpc"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/reference"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	controllerclient "sigs.k8s.io/controller-runtime/pkg/client"

	metalv1alpha1 "github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/app/metal-controller-manager/internal/api"
	"github.com/talos-systems/sidero/app/metal-controller-manager/pkg/client"
)

const (
	Port = "50100"
)

type server struct {
	api.UnimplementedAgentServer

	autoAccept bool

	c         controllerclient.Client
	clientset *kubernetes.Clientset

	metalScheme *runtime.Scheme
	recorder    record.EventRecorder
}

// CreateServer implements api.AgentServer.
func (s *server) CreateServer(ctx context.Context, in *api.CreateServerRequest) (*api.CreateServerResponse, error) {
	obj := &metalv1alpha1.Server{}

	if err := s.c.Get(ctx, types.NamespacedName{Name: in.GetSystemInformation().GetUuid()}, obj); err != nil {
		if !apierrors.IsNotFound(err) {
			return nil, err
		}

		obj = &metalv1alpha1.Server{
			TypeMeta: v1.TypeMeta{
				Kind:       "Server",
				APIVersion: metalv1alpha1.GroupVersion.Version,
			},
			ObjectMeta: v1.ObjectMeta{
				Name: in.GetSystemInformation().GetUuid(),
			},
			Spec: metalv1alpha1.ServerSpec{
				Hostname: in.GetHostname(),
				SystemInformation: &metalv1alpha1.SystemInformation{
					Manufacturer: in.GetSystemInformation().GetManufacturer(),
					ProductName:  in.GetSystemInformation().GetProductName(),
					Version:      in.GetSystemInformation().GetVersion(),
					SerialNumber: in.GetSystemInformation().GetSerialNumber(),
					SKUNumber:    in.GetSystemInformation().GetSkuNumber(),
					Family:       in.GetSystemInformation().GetFamily(),
				},
				CPU: &metalv1alpha1.CPUInformation{
					Manufacturer: in.GetCpu().GetManufacturer(),
					Version:      in.GetCpu().GetVersion(),
				},
				Accepted: s.autoAccept,
			},
		}

		if err = s.c.Create(ctx, obj); err != nil {
			return nil, err
		}

		ref, err := reference.GetReference(s.metalScheme, obj)
		if err != nil {
			return nil, err
		}

		s.recorder.Event(ref, corev1.EventTypeNormal, "Server Registration", "Server auto-registered via API.")

		log.Printf("Added %s", in.GetSystemInformation().GetUuid())
	}

	resp := &api.CreateServerResponse{}

	// Only return a wipe directive is the server is not clean *AND* it has been accepted.
	// This avoids the possibility of a random device PXE booting against us, registering, then getting blown away.
	if !obj.Status.IsClean && obj.Spec.Accepted {
		log.Printf("Server %q needs wipe", obj.Name)

		resp.Wipe = true
	}

	return resp, nil
}

// MarkServerAsWiped implements api.AgentServer.
func (s *server) MarkServerAsWiped(ctx context.Context, in *api.MarkServerAsWipedRequest) (*api.MarkServerAsWipedResponse, error) {
	obj := &metalv1alpha1.Server{}

	if err := s.c.Get(ctx, types.NamespacedName{Name: in.GetUuid()}, obj); err != nil {
		return nil, err
	}

	patchHelper, err := patch.NewHelper(obj, s.c)
	if err != nil {
		return nil, err
	}

	obj.Status.IsClean = true

	conditions.MarkTrue(obj, metalv1alpha1.ConditionPowerCycle)

	if err := patchHelper.Patch(ctx, obj); err != nil {
		return nil, err
	}

	ref, err := reference.GetReference(s.metalScheme, obj)
	if err != nil {
		return nil, err
	}

	s.recorder.Event(ref, corev1.EventTypeNormal, "Server Wipe", "Server wiped via agent.")

	resp := &api.MarkServerAsWipedResponse{}

	return resp, nil
}

// ReconcileServerAddresses implements api.AgentServer.
func (s *server) ReconcileServerAddresses(ctx context.Context, in *api.ReconcileServerAddressesRequest) (*api.ReconcileServerAddressesResponse, error) {
	obj := &metalv1alpha1.Server{}

	if err := s.c.Get(ctx, types.NamespacedName{Name: in.GetUuid()}, obj); err != nil {
		return nil, err
	}

	patchHelper, err := patch.NewHelper(obj, s.c)
	if err != nil {
		return nil, err
	}

	if func() bool {
		oldIPs := make(map[corev1.NodeAddressType]map[string]struct{})
		for _, addr := range obj.Status.Addresses {
			if oldIPs[addr.Type] == nil {
				oldIPs[addr.Type] = make(map[string]struct{})
			}
			oldIPs[addr.Type][addr.Address] = struct{}{}
		}

		newIPs := make(map[corev1.NodeAddressType]map[string]struct{})
		for _, addr := range in.GetAddress() {
			if newIPs[corev1.NodeAddressType(addr.GetType())] == nil {
				newIPs[corev1.NodeAddressType(addr.GetType())] = make(map[string]struct{})
			}
			newIPs[corev1.NodeAddressType(addr.GetType())][addr.GetAddress()] = struct{}{}
		}

		changed := false

		for typ := range newIPs {
			if oldIPs[typ] == nil {
				changed = true
				break
			}

			if !reflect.DeepEqual(newIPs[typ], oldIPs[typ]) {
				changed = true
				break
			}
		}

		if !changed {
			return false
		}

		obj.Status.Addresses = nil

		// overwrite any types submitted from the agent
		for typ := range newIPs {
			for addr := range newIPs[typ] {
				obj.Status.Addresses = append(obj.Status.Addresses, corev1.NodeAddress{
					Type:    typ,
					Address: addr,
				})
			}
		}

		// keep any types which were not sent by the agent
		for typ := range oldIPs {
			if newIPs[typ] != nil {
				continue
			}

			for addr := range oldIPs[typ] {
				obj.Status.Addresses = append(obj.Status.Addresses, corev1.NodeAddress{
					Type:    typ,
					Address: addr,
				})
			}
		}

		return true
	}() {
		if err := patchHelper.Patch(ctx, obj); err != nil {
			return nil, err
		}
	}

	resp := &api.ReconcileServerAddressesResponse{}

	return resp, nil
}

func Serve(autoAccept bool) error {
	lis, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	c, err := client.NewClient(config)
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartRecordingToSink(
		&typedcorev1.EventSinkImpl{
			Interface: clientset.CoreV1().Events(""),
		})

	recorder := eventBroadcaster.NewRecorder(
		scheme.Scheme,
		corev1.EventSource{Component: "sidero-server"})

	scheme := runtime.NewScheme()

	if err := metalv1alpha1.AddToScheme(scheme); err != nil {
		return err
	}

	api.RegisterAgentServer(s, &server{
		autoAccept:  autoAccept,
		c:           c,
		clientset:   clientset,
		recorder:    recorder,
		metalScheme: scheme,
	})

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
