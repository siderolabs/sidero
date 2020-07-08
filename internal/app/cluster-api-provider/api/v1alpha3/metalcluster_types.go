// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capiv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
)

const (
	// ClusterFinalizer allows ReconcileMetalCluster to clean up resources before removing it from the apiserver.
	ClusterFinalizer = "metalcluster.infrastructure.cluster.x-k8s.io"
)

// MetalClusterSpec defines the desired state of MetalCluster
type MetalClusterSpec struct {
	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint capiv1.APIEndpoint `json:"controlPlaneEndpoint"`
}

// MetalClusterStatus defines the observed state of MetalCluster
type MetalClusterStatus struct {
	Ready bool `json:"ready"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=metalclusters,scope=Namespaced,categories=cluster-api
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this MetalCluster belongs"
// +kubebuilder:printcolumn:name="Ready",type="boolean",JSONPath=".status.ready"
// +kubebuilder:printcolumn:name="Endpoint",type="string",priority=1,JSONPath=".spec.controlPlaneEndpoint.host",description="Control Plane Endpoint"
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

// MetalCluster is the Schema for the metalclusters API
type MetalCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MetalClusterSpec   `json:"spec,omitempty"`
	Status MetalClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MetalClusterList contains a list of MetalCluster
type MetalClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MetalCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MetalCluster{}, &MetalClusterList{})
}
