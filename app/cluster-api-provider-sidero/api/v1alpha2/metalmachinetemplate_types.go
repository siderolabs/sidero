// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MetalMachineTemplateSpec defines the desired state of MetalMachineTemplate.
type MetalMachineTemplateSpec struct {
	Template MetalMachineTemplateResource `json:"template"`
}

// MetalMachineTemplateStatus defines the observed state of MetalMachineTemplate.
type MetalMachineTemplateStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=metalmachinetemplates,scope=Namespaced,categories=cluster-api

// MetalMachineTemplate is the Schema for the metalmachinetemplates API.
type MetalMachineTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MetalMachineTemplateSpec   `json:"spec,omitempty"`
	Status MetalMachineTemplateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MetalMachineTemplateList contains a list of MetalMachineTemplate.
type MetalMachineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MetalMachineTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MetalMachineTemplate{}, &MetalMachineTemplateList{})
}

type MetalMachineTemplateResource struct {
	// Spec is the specification of the desired behavior of the machine.
	Spec MetalMachineSpec `json:"spec"`
}
