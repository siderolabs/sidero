// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha3

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServerBindingMetalMachineRefField is a reference to a field matching server binding to a metal machine.
const ServerBindingMetalMachineRefField = "spec.metalMachineRef.name"

// ServerBindingSpec defines the spec of the ServerBinding object.
type ServerBindingSpec struct {
	ServerClassRef  *corev1.ObjectReference `json:"serverClassRef,omitempty"`
	MetalMachineRef corev1.ObjectReference  `json:"metalMachineRef"`
}

// ServerBindingState defines the observed state of ServerBinding.
type ServerBindingState struct {
	// Ready is true when matching server is found.
	// +optional
	Ready bool `json:"ready"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="ServerBinding ready status"
// +kubebuilder:printcolumn:name="Server",type="string",priority=1,JSONPath=".metadata.name",description="Server ID"
// +kubebuilder:printcolumn:name="ServerClass",type="string",priority=1,JSONPath=".spec.serverClassRef.name",description="Server Class"
// +kubebuilder:printcolumn:name="MetalMachine",type="string",priority=1,JSONPath=".spec.metalMachineRef.name",description="Metal Machine"
// +kubebuilder:printcolumn:name="Cluster",type="string",priority=1,JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this ServerBinding belongs"
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

// ServerBinding defines the binding between the MetalMachine and the Server.
//
// ServerBinding always has matching ID with the Server object.
// ServerBinding optionally binds to the ServerClass which Server was picked from.
type ServerBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServerBindingSpec  `json:"spec,omitempty"`
	Status ServerBindingState `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ServerBindingList contains a list of ServerBinding.
type ServerBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServerBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ServerBinding{}, &ServerBindingList{})
}
