// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package server

import (
	"context"
	"log"
	"reflect"
	"time"

	"google.golang.org/grpc"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/reference"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	controllerclient "sigs.k8s.io/controller-runtime/pkg/client"

	metalv1alpha1 "github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/app/metal-controller-manager/internal/api"
)

type server struct {
	api.UnimplementedAgentServer

	autoAccept   bool
	insecureWipe bool
	autoBMC      bool

	c             controllerclient.Client
	scheme        *runtime.Scheme
	recorder      record.EventRecorder
	rebootTimeout time.Duration
}

// CreateServer implements api.AgentServer.
func (s *server) CreateServer(ctx context.Context, in *api.CreateServerRequest) (*api.CreateServerResponse, error) {
	obj := &metalv1alpha1.Server{}

	if err := s.c.Get(ctx, types.NamespacedName{Name: in.GetSystemInformation().GetUuid()}, obj); err != nil {
		if !apierrors.IsNotFound(err) {
			return nil, err
		}

		obj = &metalv1alpha1.Server{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Server",
				APIVersion: metalv1alpha1.GroupVersion.Version,
			},
			ObjectMeta: metav1.ObjectMeta{
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

		ref, err := reference.GetReference(s.scheme, obj)
		if err != nil {
			return nil, err
		}

		s.recorder.Event(ref, corev1.EventTypeNormal, "Server Registration", "Server auto-registered via API.")

		log.Printf("Added %s", in.GetSystemInformation().GetUuid())
	}

	resp := &api.CreateServerResponse{}

	// Make BMC and wiping decisions only if server is accepted
	// to avoid hijacking random devices that PXE boot against us.
	if obj.Spec.Accepted {
		// Respond to agent whether it should attempt bmc setup
		// We will only tell it to attempt autoconfig if there's not already data there.
		if obj.Spec.BMC == nil && s.autoBMC {
			log.Printf("Server %q needs BMC setup", obj.Name)

			resp.SetupBmc = true
		}

		// Only return a wipe directive is the server is not clean *AND* it has been accepted.
		// This avoids the possibility of a random device PXE booting against us, registering, then getting blown away.
		if !obj.Status.IsClean {
			log.Printf("Server %q needs wipe", obj.Name)

			resp.Wipe = true
			resp.InsecureWipe = s.insecureWipe
			resp.RebootTimeout = s.rebootTimeout.Seconds()
		}
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

	if err := patchHelper.Patch(ctx, obj, patch.WithOwnedConditions{
		Conditions: []clusterv1.ConditionType{metalv1alpha1.ConditionPowerCycle},
	}); err != nil {
		return nil, err
	}

	ref, err := reference.GetReference(s.scheme, obj)
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

// Heartbeat implements api.AgentServer.
func (s *server) Heartbeat(ctx context.Context, in *api.HeartbeatRequest) (*api.HeartbeatResponse, error) {
	obj := &metalv1alpha1.Server{}

	if err := s.c.Get(ctx, types.NamespacedName{Name: in.GetUuid()}, obj); err != nil {
		return nil, err
	}

	patchHelper, err := patch.NewHelper(obj, s.c)
	if err != nil {
		return nil, err
	}

	// remove the condition in case it was already set to make sure LastTransitionTime will be updated
	conditions.Delete(obj, metalv1alpha1.ConditionPowerCycle)
	conditions.MarkFalse(obj, metalv1alpha1.ConditionPowerCycle, "InProgress", clusterv1.ConditionSeverityInfo, "Server wipe in progress.")

	if err := patchHelper.Patch(ctx, obj, patch.WithOwnedConditions{
		Conditions: []clusterv1.ConditionType{metalv1alpha1.ConditionPowerCycle},
	}); err != nil {
		return nil, err
	}

	resp := &api.HeartbeatResponse{}

	return resp, nil
}

func (s *server) UpdateBMCInfo(ctx context.Context, in *api.UpdateBMCInfoRequest) (*api.UpdateBMCInfoResponse, error) {
	if in.GetBmcInfo() != nil {
		bmcSecretName := in.GetUuid() + "-bmc"

		// Create or update creds secret
		credsSecret := &corev1.Secret{}
		exists := true

		// Fetch corresponding server
		obj := &metalv1alpha1.Server{}

		if err := s.c.Get(ctx, types.NamespacedName{Name: in.GetUuid()}, obj); err != nil {
			return nil, err
		}

		// For auto-created BMC info, we will *always* drop the creds in default namespace
		// This ensures they'll come along for the ride in a "default" clusterctl move based on our cluster templates.
		if err := s.c.Get(ctx, types.NamespacedName{Namespace: corev1.NamespaceDefault, Name: bmcSecretName}, credsSecret); err != nil {
			if !apierrors.IsNotFound(err) {
				return nil, err
			}

			log.Printf("BMC secret doesn't exist for server %q, creating.", in.GetUuid())

			credsSecret.ObjectMeta = metav1.ObjectMeta{
				Namespace: corev1.NamespaceDefault,
				Name:      bmcSecretName,
				OwnerReferences: []metav1.OwnerReference{
					*metav1.NewControllerRef(obj, metalv1alpha1.GroupVersion.WithKind("Server")),
				},
			}

			exists = false
		}

		credsSecret.Data = map[string][]byte{
			"user": []byte(in.GetBmcInfo().GetUser()),
			"pass": []byte(in.GetBmcInfo().GetPass()),
		}

		if exists {
			if err := s.c.Update(ctx, credsSecret); err != nil {
				return nil, err
			}
		} else {
			if err := s.c.Create(ctx, credsSecret); err != nil {
				return nil, err
			}
		}

		// Update server spec with pointers to endpoint and creds secret

		obj.Spec.BMC = &metalv1alpha1.BMC{
			Endpoint: in.GetBmcInfo().GetIp(),
			UserFrom: &metalv1alpha1.CredentialSource{
				SecretKeyRef: &metalv1alpha1.SecretKeyRef{
					Namespace: corev1.NamespaceDefault,
					Name:      bmcSecretName,
					Key:       "user",
				},
			},
			PassFrom: &metalv1alpha1.CredentialSource{
				SecretKeyRef: &metalv1alpha1.SecretKeyRef{
					Namespace: corev1.NamespaceDefault,
					Name:      bmcSecretName,
					Key:       "pass",
				},
			},
		}

		log.Printf("Updating server %q with BMC info", in.GetUuid())

		if err := s.c.Update(ctx, obj); err != nil {
			return nil, err
		}

		ref, err := reference.GetReference(s.scheme, obj)
		if err != nil {
			return nil, err
		}

		s.recorder.Event(ref, corev1.EventTypeNormal, "BMC Update", "BMC info updated via API.")
	}

	resp := &api.UpdateBMCInfoResponse{}

	return resp, nil
}

func CreateServer(c controllerclient.Client, recorder record.EventRecorder, scheme *runtime.Scheme, autoAccept, insecureWipe, autoBMC bool, rebootTimeout time.Duration) *grpc.Server {
	s := grpc.NewServer()

	api.RegisterAgentServer(s, &server{
		autoAccept:    autoAccept,
		insecureWipe:  insecureWipe,
		autoBMC:       autoBMC,
		c:             c,
		scheme:        scheme,
		recorder:      recorder,
		rebootTimeout: rebootTimeout,
	})

	return s
}
