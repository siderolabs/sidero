// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha2

import (
	"context"
	"fmt"
	"reflect"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	siderotypes "github.com/talos-systems/sidero/app/sidero-controller-manager/pkg/types"
)

// BMC defines data about how to talk to the node via ipmitool.
type BMC struct {
	// BMC endpoint.
	Endpoint string `json:"endpoint"`
	// BMC port. Defaults to 623.
	// +optional
	Port uint32 `json:"port,omitempty"`
	// BMC user value.
	// +optional
	User string `json:"user,omitempty"`
	// Source for the user value. Cannot be used if User is not empty.
	// +optional
	UserFrom *CredentialSource `json:"userFrom,omitempty"`
	// BMC password value.
	// +optional
	Pass string `json:"pass,omitempty"`
	// Source for the password value. Cannot be used if Pass is not empty.
	// +optional
	PassFrom *CredentialSource `json:"passFrom,omitempty"`
	// BMC Interface Type. Defaults to lanplus.
	// +optional
	Interface string `json:"interface,omitempty"`
}

// CredentialSource defines a reference to the credential value.
type CredentialSource struct {
	SecretKeyRef *SecretKeyRef `json:"secretKeyRef,omitempty"`
}

// SecretKeyRef defines a ref to a given key within a secret.
type SecretKeyRef struct {
	// Namespace and name of credential secret
	// nb: can't use namespacedname here b/c it doesn't have json tags in the struct :(
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	// Key to select
	Key string `json:"key"`
}

// Resolve the value using the references.
func (source *CredentialSource) Resolve(ctx context.Context, reader client.Client) (string, error) {
	if source == nil {
		return "", nil
	}

	if source.SecretKeyRef == nil {
		return "", fmt.Errorf("missing secretKeyRef")
	}

	var secrets corev1.Secret

	if err := reader.Get(
		ctx,
		types.NamespacedName{
			Namespace: source.SecretKeyRef.Namespace,
			Name:      source.SecretKeyRef.Name,
		},
		&secrets,
	); err != nil {
		return "", fmt.Errorf("error getting secret %q: %w", source.SecretKeyRef.Name, err)
	}

	rawValue, ok := secrets.Data[source.SecretKeyRef.Key]
	if !ok {
		return "", fmt.Errorf("secret key %q is missing in secret %q", source.SecretKeyRef.Key, source.SecretKeyRef.Name)
	}

	return string(rawValue), nil
}

// ManagementAPI defines data about how to talk to the node via simple HTTP API.
type ManagementAPI struct {
	Endpoint string `json:"endpoint"`
}

type SystemInformation struct {
	Uuid         string `json:"uuid,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`
	ProductName  string `json:"productName,omitempty"`
	Version      string `json:"version,omitempty"`
	SerialNumber string `json:"serialNumber,omitempty"`
	SKUNumber    string `json:"skuNumber,omitempty"`
	Family       string `json:"family,omitempty"`
}

type Processor struct {
	Manufacturer string `json:"manufacturer,omitempty"`
	ProductName  string `json:"productName,omitempty"`
	SerialNumber string `json:"serialNumber,omitempty"`
	// Speed is in megahertz (Mhz)
	Speed       uint32 `json:"speed,omitempty"`
	CoreCount   uint32 `json:"coreCount,omitempty"`
	ThreadCount uint32 `json:"threadCount,omitempty"`
}

type ComputeInformation struct {
	TotalCoreCount   uint32       `json:"totalCoreCount,omitempty"`
	TotalThreadCount uint32       `json:"totalThreadCount,omitempty"`
	ProcessorCount   uint32       `json:"processorCount,omitempty"`
	Processors       []*Processor `json:"processors,omitempty"`
}

type MemoryModule struct {
	Manufacturer string `json:"manufacturer,omitempty"`
	ProductName  string `json:"productName,omitempty"`
	SerialNumber string `json:"serialNumber,omitempty"`
	Type         string `json:"type,omitempty"`
	// Size is in megabytes (MB)
	Size uint32 `json:"size,omitempty"`
	// Speed is in megatransfers per second (MT/S)
	Speed uint32 `json:"speed,omitempty"`
}

type MemoryInformation struct {
	TotalSize   string          `json:"totalSize,omitempty"`
	ModuleCount uint32          `json:"moduleCount,omitempty"`
	Modules     []*MemoryModule `json:"modules,omitempty"`
}

type StorageDevice struct {
	Type string `json:"type,omitempty"`
	// Size is in bytes
	Size       uint64 `json:"size,omitempty"`
	Model      string `json:"productName,omitempty"`
	Serial     string `json:"serialNumber,omitempty"`
	Name       string `json:"name,omitempty"`
	DeviceName string `json:"deviceName,omitempty"`
	UUID       string `json:"uuid,omitempty"`
	WWID       string `json:"wwid,omitempty"`
}

type StorageInformation struct {
	TotalSize   string           `json:"totalSize,omitempty"`
	DeviceCount uint32           `json:"deviceCount,omitempty"`
	Devices     []*StorageDevice `json:"devices,omitempty"`
}

type NetworkInterface struct {
	Index     uint32   `json:"index,omitempty"`
	Name      string   `json:"name,omitempty"`
	Flags     string   `json:"flags,omitempty"`
	MTU       uint32   `json:"mtu,omitempty"`
	MAC       string   `json:"mac,omitempty"`
	Addresses []string `json:"addresses,omitempty"`
}

type NetworkInformation struct {
	InterfaceCount uint32              `json:"interfaceCount,omitempty"`
	Interfaces     []*NetworkInterface `json:"interfaces,omitempty"`
}

type HardwareInformation struct {
	System  *SystemInformation  `json:"system,omitempty"`
	Compute *ComputeInformation `json:"compute,omitempty"`
	Memory  *MemoryInformation  `json:"memory,omitempty"`
	Storage *StorageInformation `json:"storage,omitempty"`
	Network *NetworkInformation `json:"network,omitempty"`
}

func (a *HardwareInformation) PartialEqual(b *HardwareInformation) bool {
	return PartialEqual(a, b)
}

func PartialEqual(a, b interface{}) bool {
	old := reflect.ValueOf(a)
	new := reflect.ValueOf(b)

	if old.Kind() == reflect.Ptr {
		old = old.Elem()
	}

	if new.Kind() == reflect.Ptr {
		new = new.Elem()
	}

	// Skip invalid or zero values, since that indicates that the user
	// did not supply the field, and does not want to compare it.
	if !old.IsValid() || old.IsZero() {
		return true
	}

	switch {
	case old.Kind() == reflect.Struct && new.Kind() == reflect.Struct:
		// Recursively compare structs
		for i := 0; i < old.NumField(); i++ {
			f1 := old.Field(i).Interface()
			f2 := new.Field(i).Interface()

			if !PartialEqual(f1, f2) {
				return false
			}
		}
	case old.Kind() == reflect.Slice && new.Kind() == reflect.Slice:
		// Skip slices where the requested slice is larger than the actual slice,
		// as that indicates that the user wants to filter for more
		// processors/memory modules/storage devices than are present.
		if old.Len() > new.Len() {
			return false
		}
		// Recursively compare slices
		for i := 0; i < old.Len(); i++ {
			f1 := old.Index(i).Interface()
			f2 := new.Index(i).Interface()

			if !PartialEqual(f1, f2) {
				return false
			}
		}
	default:
		// Directly compare values, but only if the actual value is valid.
		return new.IsValid() && old.Interface() == new.Interface()
	}

	return true
}

// ServerSpec defines the desired state of Server.
type ServerSpec struct {
	EnvironmentRef *corev1.ObjectReference `json:"environmentRef,omitempty"`
	Hardware       *HardwareInformation    `json:"hardware,omitempty"`
	Hostname       string                  `json:"hostname,omitempty"`
	BMC            *BMC                    `json:"bmc,omitempty"`
	ManagementAPI  *ManagementAPI          `json:"managementApi,omitempty"`
	ConfigPatches  []ConfigPatches         `json:"configPatches,omitempty"`
	Accepted       bool                    `json:"accepted"`
	Cordoned       bool                    `json:"cordoned,omitempty"`
	PXEBootAlways  bool                    `json:"pxeBootAlways,omitempty"`
	// BootFromDiskMethod specifies the method to exit iPXE to force boot from disk.
	//
	// If not set, controller default is used.
	// Valid values: ipxe-exit, http-404, ipxe-sanboot.
	//
	// +optional
	BootFromDiskMethod siderotypes.BootFromDisk `json:"bootFromDiskMethod,omitempty"`
	// PXEMode specifies the method to trigger PXE boot via IPMI.
	//
	// If not set, controller default is used.
	// Valid values: uefi, bios.
	//
	// +optional
	PXEMode siderotypes.PXEMode `json:"pxeMode,omitempty"`
}

const (
	// ConditionPowerCycle is used to control the powercycle flow.
	ConditionPowerCycle clusterv1.ConditionType = "PowerCycle"
	// ConditionPXEBooted is used to record the fact that server got PXE booted.
	ConditionPXEBooted clusterv1.ConditionType = "PXEBooted"
)

// ServerStatus defines the observed state of Server.
type ServerStatus struct {
	// Ready is true when server is accepted and in use.
	// +optional
	Ready bool `json:"ready"`

	// InUse is true when server is assigned to some MetalMachine.
	// +optional
	InUse bool `json:"inUse"`

	// IsClean is true when server disks are wiped.
	// +optional
	IsClean bool `json:"isClean"`

	// Conditions defines current service state of the Server.
	Conditions []clusterv1.Condition `json:"conditions,omitempty"`

	// Addresses lists discovered node IPs.
	Addresses []corev1.NodeAddress `json:"addresses,omitempty"`

	// Power is the current power state of the server: "on", "off" or "unknown".
	Power string `json:"power,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Hostname",type="string",JSONPath=".spec.hostname",description="server hostname"
// +kubebuilder:printcolumn:name="BMC IP",type="string",priority=1,JSONPath=".spec.bmc.endpoint",description="BMC IP"
// +kubebuilder:printcolumn:name="Accepted",type="boolean",JSONPath=".spec.accepted",description="indicates if the server is accepted"
// +kubebuilder:printcolumn:name="Cordoned",type="boolean",JSONPath=".spec.cordoned",description="indicates if the server is cordoned"
// +kubebuilder:printcolumn:name="Allocated",type="boolean",JSONPath=".status.inUse",description="indicates that the server has been allocated"
// +kubebuilder:printcolumn:name="Clean",type="boolean",JSONPath=".status.isClean",description="indicates if the server is clean or not"
// +kubebuilder:printcolumn:name="Power",type="string",JSONPath=".status.power",description="display the current power status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="The age of this resource"
// +kubebuilder:storageversion

// Server is the Schema for the servers API.
type Server struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServerSpec   `json:"spec,omitempty"`
	Status ServerStatus `json:"status,omitempty"`
}

func (s *Server) GetConditions() clusterv1.Conditions {
	return s.Status.Conditions
}

func (s *Server) SetConditions(conditions clusterv1.Conditions) {
	s.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// ServerList contains a list of Server.
type ServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Server `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Server{}, &ServerList{})
}
