// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha3

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/errors" //nolint:staticcheck
)

const (
	// MachineFinalizer allows ReconcileMetalMachine to clean up resources before removing it from the apiserver.
	MachineFinalizer = "metalmachine.infrastructure.cluster.x-k8s.io"
	// MetalMachineServerRefField is used to index MetalMachines on server ref it is bound to.
	MetalMachineServerRefField = "spec.serverRef.name"
)

// MetalMachineSpec defines the desired state of MetalMachine.
type MetalMachineSpec struct {
	// ProviderID is the unique identifier as specified by the cloud provider.
	// +optional
	ProviderID *string `json:"providerID,omitempty"`

	ServerRef      *corev1.ObjectReference `json:"serverRef,omitempty"`
	ServerClassRef *corev1.ObjectReference `json:"serverClassRef,omitempty"`
}

// MetalMachineStatus defines the observed state of MetalMachine.
type MetalMachineStatus struct {
	// +optional
	Ready bool `json:"ready,omitempty"`

	// Addresses contains the Metal machine associated addresses.
	Addresses []clusterv1.MachineAddress `json:"addresses,omitempty"`

	// FailureReason will be set in the event that there is a terminal problem
	// reconciling the Machine and will contain a succinct value suitable
	// for machine interpretation.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the Machine's spec or the configuration of
	// the controller, and that manual intervention is required. Examples
	// of terminal errors would be invalid combinations of settings in the
	// spec, values that are unsupported by the controller, or the
	// responsible controller itself being critically misconfigured.
	//
	// Any transient errors that occur during the reconciliation of Machines
	// can be added as events to the Machine object and/or logged in the
	// controller's output.
	// +optional
	FailureReason *errors.MachineStatusError `json:"failureReason,omitempty"`

	// FailureMessage will be set in the event that there is a terminal problem
	// reconciling the Machine and will contain a more verbose string suitable
	// for logging and human consumption.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the Machine's spec or the configuration of
	// the controller, and that manual intervention is required. Examples
	// of terminal errors would be invalid combinations of settings in the
	// spec, values that are unsupported by the controller, or the
	// responsible controller itself being critically misconfigured.
	//
	// Any transient errors that occur during the reconciliation of Machines
	// can be added as events to the Machine object and/or logged in the
	// controller's output.
	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"`

	// Conditions defines current state of the MetalMachine.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=metalmachines,scope=Namespaced,categories=cluster-api
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="MetalMachine ready status"
// +kubebuilder:printcolumn:name="Cluster",type="string",priority=1,JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this MetalMachine belongs"
// +kubebuilder:printcolumn:name="Machine",type="string",priority=1,JSONPath=".metadata.ownerReferences[?(@.kind==\"Machine\")].name",description="Machine object to which this MetalMachine belongs"
// +kubebuilder:printcolumn:name="Server",type="string",priority=1,JSONPath=".spec.serverRef.name",description="Server ID"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="The age of this resource"
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

// MetalMachine is the Schema for the metalmachines API.
type MetalMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MetalMachineSpec   `json:"spec,omitempty"`
	Status MetalMachineStatus `json:"status,omitempty"`
}

// GetConditions returns the set of conditions for this object.
func (in *MetalMachine) GetConditions() clusterv1.Conditions {
	return in.Status.Conditions
}

// SetConditions sets the conditions on this object.
func (in *MetalMachine) SetConditions(conditions clusterv1.Conditions) {
	in.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// MetalMachineList contains a list of MetalMachine.
type MetalMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MetalMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MetalMachine{}, &MetalMachineList{})
}
