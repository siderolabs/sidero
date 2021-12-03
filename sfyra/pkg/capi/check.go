// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package capi

import (
	"context"
	"fmt"

	cacpt "github.com/talos-systems/cluster-api-control-plane-provider-talos/api/v1alpha3"
	"github.com/talos-systems/go-retry/retry"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	capiv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CheckClusterReady verifies that cluster ready from the CAPI point of view.
func CheckClusterReady(ctx context.Context, metalClient client.Client, clusterName string) error {
	var cluster capiv1.Cluster

	if err := metalClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: clusterName}, &cluster); err != nil {
		return err
	}

	ready := false

	for _, cond := range cluster.Status.Conditions {
		if cond.Type == capiv1.ReadyCondition && cond.Status == corev1.ConditionTrue {
			ready = true

			break
		}
	}

	if !ready {
		return retry.ExpectedError(fmt.Errorf("cluster is not ready"))
	}

	var controlPlane cacpt.TalosControlPlane

	if err := metalClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: clusterName + "-cp"}, &controlPlane); err != nil {
		return err
	}

	if !controlPlane.Status.Ready {
		return retry.ExpectedError(fmt.Errorf("control plane is not ready"))
	}

	if !controlPlane.Status.Initialized {
		return retry.ExpectedError(fmt.Errorf("control plane is not initialized"))
	}

	if controlPlane.Status.Replicas != controlPlane.Status.ReadyReplicas {
		return retry.ExpectedError(fmt.Errorf("control plane replicas %d != ready replicas %d", controlPlane.Status.Replicas, controlPlane.Status.ReadyReplicas))
	}

	var machineDeployment capiv1.MachineDeployment

	if err := metalClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: clusterName + "-workers"}, &machineDeployment); err != nil {
		return err
	}

	if machineDeployment.Status.GetTypedPhase() != capiv1.MachineDeploymentPhaseRunning {
		return retry.ExpectedError(fmt.Errorf("machineDeployment phase is %s", machineDeployment.Status.GetTypedPhase()))
	}

	if machineDeployment.Status.Replicas != machineDeployment.Status.ReadyReplicas {
		return retry.ExpectedError(fmt.Errorf("machineDeployment replicas %d != ready replicas %d", machineDeployment.Status.Replicas, machineDeployment.Status.ReadyReplicas))
	}

	return nil
}
