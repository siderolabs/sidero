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

const (
	// TalosConfigValidatedCondition reports when talos has loaded and validated the config
	// for the machine.
	TalosConfigValidatedCondition clusterv1.ConditionType = "TalosConfigValidated"

	// TalosConfigValidationFailedReason (Severity=Error) documents that Talos config validation has failed.
	TalosConfigValidationFailedReason = "TalosConfigValidationFailed"
)

const (
	// TalosConfigLoadedCondition reports when talos has loaded the config
	// for the machine.
	TalosConfigLoadedCondition clusterv1.ConditionType = "TalosConfigLoaded"

	// TalosConfigLoadedationFailedReason (Severity=Error) documents that Talos config validation has failed.
	TalosConfigLoadFailedReason = "TalosConfigLoadFailed"
)

const (
	// TalosInstalledCondition reports when Talos OS was successfully installed on the node.
	TalosInstalledCondition clusterv1.ConditionType = "TalosInstalled"

	// TalosInstallationInProgressReason (Severity=Info) documents that Talos installation is in progress.
	TalosInstallationInProgressReason = "TalosInstallationInProgress"

	// TalosInstallationFailedReason (Severity=Error) documents that Talos installer has failed.
	TalosInstallationFailedReason = "TalosInstallationFailed"
)
