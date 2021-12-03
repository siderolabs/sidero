// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +kubebuilder:object:generate=true
// +groupName=controlplane.cluster.x-k8s.io
package v1alpha3

import clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

// Conditions and condition Reasons for the TalosControlPlane object

const (
	// ProviderSetCondition reports when the nodes have .spec.Provider ID field set.
	ProviderSetCondition clusterv1.ConditionType = "ProviderSet"

	// ProviderUpdateFailedReason (Severity=Warning) documents that controller failed
	// to set ProviderID labels on all nodes.
	ProviderUpdateFailedReason = "ProviderUpdateFailed"
)
