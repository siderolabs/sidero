// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	capbt "github.com/siderolabs/cluster-api-control-plane-provider-talos/api/v1alpha3"
	"github.com/siderolabs/go-retry/retry"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	capiv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/siderolabs/sidero/sfyra/pkg/capi"
	"github.com/siderolabs/sidero/sfyra/pkg/talos"
)

type ScaleCallBack func(runtime.Object) error

// TestScaleControlPlaneUp verifies that the control plane can scale up.
func TestScaleControlPlaneUp(ctx context.Context, metalClient client.Client, capiCluster talos.Cluster) TestFunc {
	return func(t *testing.T) {
		err := scaleControlPlane(ctx, metalClient, 3)
		require.NoError(t, err)

		err = verifyClusterHealth(ctx, metalClient, capiCluster, t)
		require.NoError(t, err)
	}
}

// TestScaleControlPlaneDown verifies that the control plane can scale down.
func TestScaleControlPlaneDown(ctx context.Context, metalClient client.Client, capiCluster talos.Cluster) TestFunc {
	return func(t *testing.T) {
		err := scaleControlPlane(ctx, metalClient, 1)
		require.NoError(t, err)

		err = verifyClusterHealth(ctx, metalClient, capiCluster, t)
		require.NoError(t, err)
	}
}

// TestScaleWorkersUp verifies that the workers can scale up.
func TestScaleWorkersUp(ctx context.Context, metalClient client.Client, capiCluster talos.Cluster) TestFunc {
	return func(t *testing.T) {
		err := scaleWorkers(ctx, metalClient, 3)
		require.NoError(t, err)

		err = verifyClusterHealth(ctx, metalClient, capiCluster, t)
		require.NoError(t, err)
	}
}

// TestScaleWorkersDown verifies that the workers can scale down.
func TestScaleWorkersDown(ctx context.Context, metalClient client.Client, capiCluster talos.Cluster) TestFunc {
	return func(t *testing.T) {
		err := scaleWorkers(ctx, metalClient, 1)
		require.NoError(t, err)

		err = verifyClusterHealth(ctx, metalClient, capiCluster, t)
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
		o := obj.(*capiv1.MachineDeployment)

		if capiv1.MachineDeploymentPhase(o.Status.Phase) != capiv1.MachineDeploymentPhaseRunning {
			return fmt.Errorf("expected %s phase, got %s", capiv1.MachineDeploymentPhaseRunning, o.Status.Phase)
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
		o := obj.(*capiv1.MachineDeployment)

		o.Spec.Replicas = &replicas

		return nil
	}

	var obj capiv1.MachineDeployment

	return scale(ctx, metalClient, "management-cluster-workers", &obj, set, verify)
}

func scale(ctx context.Context, metalClient client.Client, name string, obj client.Object, set, verify ScaleCallBack) error {
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
		obj = cleanObj.DeepCopyObject().(client.Object)

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

func verifyClusterHealth(ctx context.Context, metalClient client.Reader, capiCluster talos.Cluster, t *testing.T) error {
	t.Log("verifying cluster health")

	cluster, err := capi.NewCluster(ctx, metalClient, managementClusterName, capiCluster.BridgeIP())
	if err != nil {
		return err
	}

	err = cluster.Health(ctx)
	if err != nil {
		return err
	}

	return nil
}
