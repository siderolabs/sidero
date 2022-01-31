// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/talos-systems/go-retry/retry"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	capiv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "github.com/talos-systems/sidero/app/caps-controller-manager/api/v1alpha3"
	metalv1 "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha2"
)

// TestMachineDeploymentReconcile verifies that machine deployment can reconcile delete machines.
func TestMachineDeploymentReconcile(ctx context.Context, metalClient client.Client) TestFunc {
	return func(t *testing.T) {
		var machineDeployment capiv1.MachineDeployment

		const machineDeploymentName = "management-cluster-workers"

		err := metalClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: machineDeploymentName}, &machineDeployment)
		require.NoError(t, err)

		var machines capiv1.MachineList

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
			var machineDeployment capiv1.MachineDeployment

			err = metalClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: machineDeploymentName}, &machineDeployment)
			if err != nil {
				return err
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
				return err
			}

			if capiv1.MachineDeploymentPhase(machineDeployment.Status.Phase) != capiv1.MachineDeploymentPhaseRunning {
				return retry.ExpectedError(fmt.Errorf("expected %s phase, got %s", capiv1.MachineDeploymentPhaseRunning, machineDeployment.Status.Phase))
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

// TestServerBindingReconcile verifies that server binding controller can reconcile missing ServerBindings.
func TestServerBindingReconcile(ctx context.Context, metalClient client.Client) TestFunc {
	return func(t *testing.T) {
		var serverBindingList infrav1.ServerBindingList

		require.NoError(t, metalClient.List(ctx, &serverBindingList))

		if len(serverBindingList.Items) < 1 {
			t.Fatal("no serverbindings found")
		}

		// pick any serverbinding and delete it
		serverBindingToDelete := serverBindingList.Items[0]

		require.NoError(t, metalClient.Delete(ctx, &serverBindingToDelete))

		// verify that matching server doesn't become unallocated for 1 minute
		start := time.Now()

		for time.Since(start) < time.Minute {
			var server metalv1.Server

			require.NoError(t, metalClient.Get(ctx, types.NamespacedName{Name: serverBindingToDelete.Name}, &server))

			require.True(t, server.Status.InUse)
		}

		// server binding should have been re-created
		var serverBinding infrav1.ServerBinding

		require.NoError(t, metalClient.Get(ctx, types.NamespacedName{Name: serverBindingToDelete.Name}, &serverBinding))

		assert.Equal(t, serverBinding.Spec.MetalMachineRef, serverBindingToDelete.Spec.MetalMachineRef)
		assert.Equal(t, serverBinding.Labels, serverBindingToDelete.Labels)

		if serverBindingToDelete.Spec.ServerClassRef == nil {
			assert.Nil(t, serverBinding.Spec.ServerClassRef)
		} else {
			assert.Equal(t, serverBindingToDelete.Spec.ServerClassRef.Name, serverBinding.Spec.ServerClassRef.Name)
		}
	}
}

// TestMetalMachineServerRefReconcile verifies that metal machine controller can reconcile missing MetalMachine.Spec.ServerRef.
//
// This simulates failure in two-step process of metal machine server allocation: serverbinding got created, but metalmachine's server
// ref wasn't set.
func TestMetalMachineServerRefReconcile(ctx context.Context, metalClient client.Client) TestFunc {
	return func(t *testing.T) {
		var serverBindingList infrav1.ServerBindingList

		require.NoError(t, metalClient.List(ctx, &serverBindingList))

		if len(serverBindingList.Items) < 1 {
			t.Fatal("no serverbindings found")
		}

		// pick any serverbinding
		serverBinding := serverBindingList.Items[0]

		// get matching metalmachine
		var metalMachine infrav1.MetalMachine

		require.NoError(t, metalClient.Get(ctx, types.NamespacedName{Namespace: serverBinding.Spec.MetalMachineRef.Namespace, Name: serverBinding.Spec.MetalMachineRef.Name}, &metalMachine))

		patchHelper, err := patch.NewHelper(&metalMachine, metalClient)
		require.NoError(t, err)

		metalMachine.Spec.ServerRef = nil

		// nullify server ref
		require.NoError(t, patchHelper.Patch(ctx, &metalMachine))

		require.NoError(t, retry.Constant(time.Minute, retry.WithUnits(5*time.Second)).Retry(func() error {
			if err := metalClient.Get(ctx, types.NamespacedName{Namespace: serverBinding.Spec.MetalMachineRef.Namespace, Name: serverBinding.Spec.MetalMachineRef.Name}, &metalMachine); err != nil {
				return err
			}

			if metalMachine.Spec.ServerRef == nil {
				return retry.ExpectedError(fmt.Errorf("still missing server ref"))
			}

			return nil
		}))

		assert.Equal(t, serverBinding.Name, metalMachine.Spec.ServerRef.Name)
	}
}
