// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/talos-systems/go-retry/retry"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TestMachineDeploymentReconcile verifies that machine deployment can reconcile delete machines.
func TestMachineDeploymentReconcile(ctx context.Context, metalClient client.Client) TestFunc {
	return func(t *testing.T) {
		var machineDeployment v1alpha3.MachineDeployment

		const machineDeploymentName = "management-cluster-workers"

		err := metalClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: machineDeploymentName}, &machineDeployment)
		require.NoError(t, err)

		var machines v1alpha3.MachineList

		labelSelector, err := labels.Parse(machineDeployment.Status.Selector)
		require.NoError(t, err)

		err = metalClient.List(ctx, &machines, client.MatchingLabelsSelector{Selector: labelSelector})
		require.NoError(t, err)

		replicas := int32(len(machines.Items))

		for _, machine := range machines.Items {
			machine := machine

			err = metalClient.Delete(ctx, &machine)
			require.NoError(t, err)
		}

		// first, controller should pick up the fact that some replicas are missing
		err = retry.Constant(1*time.Minute, retry.WithUnits(100*time.Millisecond)).Retry(func() error {
			var machineDeployment v1alpha3.MachineDeployment

			err = metalClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: machineDeploymentName}, &machineDeployment)
			if err != nil {
				return retry.UnexpectedError(err)
			}

			if machineDeployment.Status.UnavailableReplicas != replicas {
				return retry.ExpectedError(fmt.Errorf("expected %d unavailable replicas, got %d", replicas, machineDeployment.Status.ReadyReplicas))
			}

			return nil
		})
		require.NoError(t, err)

		// next, check that replicas get reconciled
		err = retry.Constant(10*time.Minute, retry.WithUnits(10*time.Second)).Retry(func() error {
			err = metalClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: machineDeploymentName}, &machineDeployment)
			if err != nil {
				return retry.UnexpectedError(err)
			}

			if v1alpha3.MachineDeploymentPhase(machineDeployment.Status.Phase) != v1alpha3.MachineDeploymentPhaseRunning {
				return retry.ExpectedError(fmt.Errorf("expected %s phase, got %s", v1alpha3.MachineDeploymentPhaseRunning, machineDeployment.Status.Phase))
			}

			if machineDeployment.Status.Replicas != replicas {
				return retry.ExpectedError(fmt.Errorf("expected %d replicas, got %d", replicas, machineDeployment.Status.Replicas))
			}

			if machineDeployment.Status.ReadyReplicas != replicas {
				return retry.ExpectedError(fmt.Errorf("expected %d ready replicas, got %d", replicas, machineDeployment.Status.ReadyReplicas))
			}

			return nil
		})
		require.NoError(t, err)
	}
}
