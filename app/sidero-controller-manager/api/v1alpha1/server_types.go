// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha1

import (
	"context"
	"fmt"
	"reflect"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	Manufacturer string `json:"manufacturer,omitempty"`
	ProductName  string `json:"productName,omitempty"`
	Version      string `json:"version,omitempty"`
	SerialNumber string `json:"serialNumber,omitempty"`
	SKUNumber    string `json:"skuNumber,omitempty"`
	Family       string `json:"family,omitempty"`
}

func (a *SystemInformation) PartialEqual(b *SystemInformation) bool {
	return PartialEqual(a, b)
}

type CPUInformation struct {
	Manufacturer string `json:"manufacturer,omitempty"`
	Version      string `json:"version,omitempty"`
}

func (a *CPUInformation) PartialEqual(b *CPUInformation) bool {
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

	for i := 0; i < old.NumField(); i++ {
		if old.Field(i).IsZero() {
			// Skip zero values, since that indicates that the user did not supply
			// the field, and does not want to compare it.
			continue
		}

		f1 := old.Field(i).Interface()
		f2 := new.Field(i).Interface()

		if f1 != f2 {
			return false
		}
	}

	return true
}

// ServerSpec defines the desired state of Server.
type ServerSpec struct {
	EnvironmentRef    *corev1.ObjectReference `json:"environmentRef,omitempty"`
	Hostname          string                  `json:"hostname,omitempty"`
	SystemInformation *SystemInformation      `json:"system,omitempty"`
	CPU               *CPUInformation         `json:"cpu,omitempty"`
	BMC               *BMC                    `json:"bmc,omitempty"`
	ManagementAPI     *ManagementAPI          `json:"managementApi,omitempty"`
	ConfigPatches     []ConfigPatches         `json:"configPatches,omitempty"`
	Accepted          bool                    `json:"accepted"`
	PXEBootAlways     bool                    `json:"pxeBootAlways,omitempty"`
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
// +kubebuilder:printcolumn:name="Accepted",type="boolean",JSONPath=".spec.accepted",description="indicates if the server is accepted"
// +kubebuilder:printcolumn:name="Allocated",type="boolean",JSONPath=".status.inUse",description="indicates that the server has been allocated"
// +kubebuilder:printcolumn:name="Clean",type="boolean",JSONPath=".status.isClean",description="indicates if the server is clean or not"
// +kubebuilder:printcolumn:name="Power",type="string",JSONPath=".status.power",description="display the current power status"

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
