// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package server

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/talos-systems/grpc-proxy/proxy"
	"google.golang.org/grpc"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/reference"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterctl "sigs.k8s.io/cluster-api/cmd/clusterctl/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	controllerclient "sigs.k8s.io/controller-runtime/pkg/client"

	metalv1 "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha2"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/api"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/pkg/constants"
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
	obj := &metalv1.Server{}
	uuid := in.GetHardware().GetSystem().GetUuid()

	if err := s.c.Get(ctx, types.NamespacedName{Name: uuid}, obj); err != nil {
		if !apierrors.IsNotFound(err) {
			return nil, err
		}

		obj = &metalv1.Server{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Server",
				APIVersion: metalv1.GroupVersion.Version,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: uuid,
			},
			Spec: metalv1.ServerSpec{
				Hardware: MapHardwareInformation(in.GetHardware()),
				Hostname: in.GetHostname(),
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

		log.Printf("Added %s", uuid)
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
	obj := &metalv1.Server{}

	if err := s.c.Get(ctx, types.NamespacedName{Name: in.GetUuid()}, obj); err != nil {
		return nil, err
	}

	patchHelper, err := patch.NewHelper(obj, s.c)
	if err != nil {
		return nil, err
	}

	obj.Status.IsClean = true

	conditions.MarkTrue(obj, metalv1.ConditionPowerCycle)

	if err := patchHelper.Patch(ctx, obj, patch.WithOwnedConditions{
		Conditions: []clusterv1.ConditionType{metalv1.ConditionPowerCycle},
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
	obj := &metalv1.Server{}

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
	obj := &metalv1.Server{}

	if err := s.c.Get(ctx, types.NamespacedName{Name: in.GetUuid()}, obj); err != nil {
		return nil, err
	}

	patchHelper, err := patch.NewHelper(obj, s.c)
	if err != nil {
		return nil, err
	}

	// remove the condition in case it was already set to make sure LastTransitionTime will be updated
	conditions.Delete(obj, metalv1.ConditionPowerCycle)
	conditions.MarkFalse(obj, metalv1.ConditionPowerCycle, "InProgress", clusterv1.ConditionSeverityInfo, "Server wipe in progress.")

	if err := patchHelper.Patch(ctx, obj, patch.WithOwnedConditions{
		Conditions: []clusterv1.ConditionType{metalv1.ConditionPowerCycle},
	}); err != nil {
		return nil, err
	}

	resp := &api.HeartbeatResponse{}

	return resp, nil
}

func (s *server) UpdateBMCInfo(ctx context.Context, in *api.UpdateBMCInfoRequest) (*api.UpdateBMCInfoResponse, error) {
	bmcInfo := in.GetBmcInfo()

	// Fetch corresponding server
	obj := &metalv1.Server{}

	if err := s.c.Get(ctx, types.NamespacedName{Name: in.GetUuid()}, obj); err != nil {
		return nil, err
	}

	// Create a BMC struct if non-existent
	if obj.Spec.BMC == nil {
		obj.Spec.BMC = &metalv1.BMC{}
	}

	// Update bmc info with IP if we've got it.
	if ip := in.GetBmcInfo().GetIp(); ip != "" {
		obj.Spec.BMC.Endpoint = ip
	}

	// Update bmc info with port
	obj.Spec.BMC.Port = constants.DefaultBMCPort

	port := in.GetBmcInfo().GetPort()
	if port != 0 {
		obj.Spec.BMC.Port = port
	}

	// Generate or update bmc secret if we have creds
	if bmcInfo.User != "" && bmcInfo.Pass != "" {
		// Create or update creds secret
		credsSecret := &corev1.Secret{}
		exists := true

		// For auto-created BMC info, we will add the "move" label from clusterctl's labels.
		// This ensures they'll come along for the ride in a clusterctl move.
		bmcSecretName := in.GetUuid() + "-bmc"

		if err := s.c.Get(ctx, types.NamespacedName{Namespace: corev1.NamespaceDefault, Name: bmcSecretName}, credsSecret); err != nil {
			if !apierrors.IsNotFound(err) {
				return nil, err
			}

			log.Printf("BMC secret doesn't exist for server %q, creating.", in.GetUuid())

			credsSecret.ObjectMeta = metav1.ObjectMeta{
				Namespace: corev1.NamespaceDefault,
				Name:      bmcSecretName,
				OwnerReferences: []metav1.OwnerReference{
					*metav1.NewControllerRef(obj, metalv1.GroupVersion.WithKind("Server")),
				},
				Labels: map[string]string{
					clusterctl.ClusterctlMoveLabelName: "",
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
		obj.Spec.BMC.UserFrom = &metalv1.CredentialSource{
			SecretKeyRef: &metalv1.SecretKeyRef{
				Namespace: corev1.NamespaceDefault,
				Name:      bmcSecretName,
				Key:       "user",
			},
		}

		obj.Spec.BMC.PassFrom = &metalv1.CredentialSource{
			SecretKeyRef: &metalv1.SecretKeyRef{
				Namespace: corev1.NamespaceDefault,
				Name:      bmcSecretName,
				Key:       "pass",
			},
		}
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

	resp := &api.UpdateBMCInfoResponse{}

	return resp, nil
}

func CreateServer(c controllerclient.Client, recorder record.EventRecorder, scheme *runtime.Scheme, autoAccept, insecureWipe, autoBMC bool, rebootTimeout time.Duration) *grpc.Server {
	s := grpc.NewServer(
		// proxy pass unknown requests to sub-components
		grpc.CustomCodec(proxy.Codec()), //nolint:staticcheck
		grpc.UnknownServiceHandler(
			proxy.TransparentHandler(
				director,
			)),
	)

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

func MapHardwareInformation(hw *api.HardwareInformation) *metalv1.HardwareInformation {
	processors := make([]*metalv1.Processor, hw.GetCompute().GetProcessorCount())
	for i, v := range hw.GetCompute().GetProcessors() {
		processors[i] = &metalv1.Processor{
			Manufacturer: v.GetManufacturer(),
			ProductName:  v.GetProductName(),
			SerialNumber: v.GetSerialNumber(),
			Speed:        v.GetSpeed(),
			CoreCount:    v.GetCoreCount(),
			ThreadCount:  v.GetThreadCount(),
		}
	}

	memoryModules := make([]*metalv1.MemoryModule, hw.GetMemory().GetModuleCount())
	for i, v := range hw.GetMemory().GetModules() {
		memoryModules[i] = &metalv1.MemoryModule{
			Manufacturer: v.GetManufacturer(),
			ProductName:  v.GetProductName(),
			SerialNumber: v.GetSerialNumber(),
			Type:         v.GetType(),
			Size:         v.GetSize(),
			Speed:        v.GetSpeed(),
		}
	}

	storageDevices := make([]*metalv1.StorageDevice, hw.GetStorage().GetDeviceCount())
	for i, v := range hw.GetStorage().GetDevices() {
		storageDevices[i] = &metalv1.StorageDevice{
			Type:       v.GetType().String(),
			Size:       v.GetSize(),
			Model:      v.GetModel(),
			Serial:     v.GetSerial(),
			Name:       v.GetName(),
			DeviceName: v.GetDeviceName(),
			UUID:       v.GetUuid(),
			WWID:       v.GetWwid(),
		}
	}

	networkInterfaces := make([]*metalv1.NetworkInterface, hw.GetNetwork().GetInterfaceCount())
	for i, v := range hw.GetNetwork().GetInterfaces() {
		networkInterfaces[i] = &metalv1.NetworkInterface{
			Index:     v.GetIndex(),
			Name:      v.GetName(),
			Flags:     v.GetFlags(),
			MTU:       v.GetMtu(),
			MAC:       v.GetMac(),
			Addresses: v.GetAddresses(),
		}
	}

	return &metalv1.HardwareInformation{
		System: &metalv1.SystemInformation{
			Uuid:         hw.GetSystem().GetUuid(),
			Manufacturer: hw.GetSystem().GetManufacturer(),
			ProductName:  hw.GetSystem().GetProductName(),
			Version:      hw.GetSystem().GetVersion(),
			SerialNumber: hw.GetSystem().GetSerialNumber(),
			SKUNumber:    hw.GetSystem().GetSkuNumber(),
			Family:       hw.GetSystem().GetFamily(),
		},
		Compute: &metalv1.ComputeInformation{
			TotalCoreCount:   hw.GetCompute().GetTotalCoreCount(),
			TotalThreadCount: hw.GetCompute().GetTotalThreadCount(),
			ProcessorCount:   hw.GetCompute().GetProcessorCount(),
			Processors:       processors,
		},
		Memory: &metalv1.MemoryInformation{
			TotalSize:   fmt.Sprintf("%d GB", hw.GetMemory().GetTotalSize()/1024),
			ModuleCount: hw.GetMemory().GetModuleCount(),
			Modules:     memoryModules,
		},
		Storage: &metalv1.StorageInformation{
			TotalSize:   fmt.Sprintf("%d GB", hw.GetStorage().GetTotalSize()/1024/1024/1024),
			DeviceCount: hw.GetStorage().GetDeviceCount(),
			Devices:     storageDevices,
		},
		Network: &metalv1.NetworkInformation{
			InterfaceCount: hw.GetNetwork().GetInterfaceCount(),
			Interfaces:     networkInterfaces,
		},
	}
}
