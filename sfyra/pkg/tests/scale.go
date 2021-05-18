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
	capbt "github.com/talos-systems/cluster-api-control-plane-provider-talos/api/v1alpha3"
	"github.com/talos-systems/go-retry/retry"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/talos-systems/sidero/sfyra/pkg/capi"
	"github.com/talos-systems/sidero/sfyra/pkg/vm"
)

type ScaleCallBack func(runtime.Object) error

// TestScaleControlPlaneUp verifies that the control plane can scale up.
func TestScaleControlPlaneUp(ctx context.Context, metalClient client.Client, vmSet *vm.Set) TestFunc {
	return func(t *testing.T) {
		err := scaleControlPlane(ctx, metalClient, 3)
		require.NoError(t, err)

		err = verifyClusterHealth(ctx, metalClient, vmSet, t)
		require.NoError(t, err)
	}
}

// TestScaleControlPlaneDown verifies that the control plane can scale down.
func TestScaleControlPlaneDown(ctx context.Context, metalClient client.Client, vmSet *vm.Set) TestFunc {
	return func(t *testing.T) {
		err := scaleControlPlane(ctx, metalClient, 1)
		require.NoError(t, err)

		err = verifyClusterHealth(ctx, metalClient, vmSet, t)
		require.NoError(t, err)
	}
}

// TestScaleWorkersUp verifies that the workers can scale up.
func TestScaleWorkersUp(ctx context.Context, metalClient client.Client, vmSet *vm.Set) TestFunc {
	return func(t *testing.T) {
		err := scaleWorkers(ctx, metalClient, 3)
		require.NoError(t, err)

		err = verifyClusterHealth(ctx, metalClient, vmSet, t)
		require.NoError(t, err)
	}
}

// TestScaleWorkersDown verifies that the workers can scale down.
func TestScaleWorkersDown(ctx context.Context, metalClient client.Client, vmSet *vm.Set) TestFunc {
	return func(t *testing.T) {
		err := scaleWorkers(ctx, metalClient, 1)
		require.NoError(t, err)

		err = verifyClusterHealth(ctx, metalClient, vmSet, t)
		require.NoError(t, err)
	}
}

func scaleControlPlane(ctx context.Context, metalClient client.Client, replicas int32) error {
	verify := func(obj runtime.Object) error {
		o := obj.(*capbt.TalosControlPlane)

		if o.Status.Replicas != replicas {
			return fmt.Errorf("expected %d replicas, got %d", replicas, o.Status.ReadyReplicas)
		}

		if o.Status.ReadyReplicas != replicas {
			return fmt.Errorf("expected %d ready replicas, got %d", replicas, o.Status.ReadyReplicas)
		}

		return nil
	}

	set := func(obj runtime.Object) error {
		o := obj.(*capbt.TalosControlPlane)

		o.Spec.Replicas = &replicas

		return nil
	}

	var obj capbt.TalosControlPlane

	return scale(ctx, metalClient, "management-cluster-cp", &obj, set, verify)
}

func scaleWorkers(ctx context.Context, metalClient client.Client, replicas int32) error {
	verify := func(obj runtime.Object) error {
		o := obj.(*v1alpha3.MachineDeployment)

		if v1alpha3.MachineDeploymentPhase(o.Status.Phase) != v1alpha3.MachineDeploymentPhaseRunning {
			return fmt.Errorf("expected %s phase, got %s", v1alpha3.MachineDeploymentPhaseRunning, o.Status.Phase)
		}

		if o.Status.Replicas != replicas {
			return fmt.Errorf("expected %d replicas, got %d", replicas, o.Status.Replicas)
		}

		if o.Status.ReadyReplicas != replicas {
			return fmt.Errorf("expected %d ready replicas, got %d", replicas, o.Status.ReadyReplicas)
		}

		return nil
	}

	set := func(obj runtime.Object) error {
		o := obj.(*v1alpha3.MachineDeployment)

		o.Spec.Replicas = &replicas

		return nil
	}

	var obj v1alpha3.MachineDeployment

	return scale(ctx, metalClient, "management-cluster-workers", &obj, set, verify)
}

func scale(ctx context.Context, metalClient client.Client, name string, obj runtime.Object, set, verify ScaleCallBack) error {
	cleanObj := obj.DeepCopyObject()

	err := metalClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: name}, obj)
	if err != nil {
		return err
	}

	err = set(obj)
	if err != nil {
		return err
	}

	err = metalClient.Update(ctx, obj, &client.UpdateOptions{
		FieldManager: "sfyra",
	})
	if err != nil {
		return err
	}

	err = retry.Constant(10*time.Minute, retry.WithUnits(10*time.Second)).Retry(func() error {
		obj = cleanObj.DeepCopyObject()

		err := metalClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: name}, obj)
		if err != nil {
			return err
		}

		err = verify(obj)
		if err != nil {
			return retry.ExpectedError(err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	err = retry.Constant(time.Minute, retry.WithUnits(10*time.Second)).Retry(func() error {
		return capi.CheckClusterReady(ctx, metalClient, managementClusterName)
	})
	if err != nil {
		return err
	}

	return nil
}

func verifyClusterHealth(ctx context.Context, metalClient client.Reader, vmSet *vm.Set, t *testing.T) error {
	t.Log("verifying cluster health")

	cluster, err := capi.NewCluster(ctx, metalClient, managementClusterName, vmSet.BridgeIP())
	if err != nil {
		return err
	}

	err = cluster.Health(ctx)
	if err != nil {
		return err
	}

	return nil
}
