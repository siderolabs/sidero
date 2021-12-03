// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tests

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/talos-systems/go-retry/retry"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	capiv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	capiclient "sigs.k8s.io/cluster-api/cmd/clusterctl/client"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/talos-systems/sidero/sfyra/pkg/capi"
	"github.com/talos-systems/sidero/sfyra/pkg/loadbalancer"
	"github.com/talos-systems/sidero/sfyra/pkg/talos"
	"github.com/talos-systems/sidero/sfyra/pkg/vm"
)

func deployCluster(ctx context.Context, t *testing.T, metalClient client.Client, capiCluster talos.Cluster, vmSet *vm.Set,
	capiManager *capi.Manager, clusterName, serverClassName string, loadbalancerPort int, controlPlaneNodes, workerNodes int64, talosVersion, kubernetesVersion string) (*loadbalancer.ControlPlane, *capi.Cluster) {
	t.Logf("deploying cluster %q from server class %q with loadbalancer port %d", clusterName, serverClassName, loadbalancerPort)

	kubeconfig, err := capiManager.GetKubeconfig(ctx)
	require.NoError(t, err)

	config, err := capiCluster.KubernetesClient().K8sRestConfig(ctx)
	require.NoError(t, err)

	capiClient := capiManager.GetManagerClient()

	loadbalancer, err := loadbalancer.NewControlPlane(metalClient, vmSet.BridgeIP(), loadbalancerPort, "default", clusterName, false)
	require.NoError(t, err)

	t.Setenv("CONTROL_PLANE_ENDPOINT", vmSet.BridgeIP().String())
	t.Setenv("CONTROL_PLANE_PORT", strconv.Itoa(loadbalancerPort))
	t.Setenv("CONTROL_PLANE_SERVERCLASS", serverClassName)
	t.Setenv("WORKER_SERVERCLASS", serverClassName)
	t.Setenv("KUBERNETES_VERSION", kubernetesVersion)
	t.Setenv("TALOS_VERSION", talosVersion)

	templateOptions := capiclient.GetClusterTemplateOptions{
		Kubeconfig:               kubeconfig,
		ClusterName:              clusterName,
		ControlPlaneMachineCount: &controlPlaneNodes,
		WorkerMachineCount:       &workerNodes,
	}

	template, err := capiClient.GetClusterTemplate(templateOptions)
	require.NoError(t, err)

	dc, err := discovery.NewDiscoveryClientForConfig(config)
	require.NoError(t, err)

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	dyn, err := dynamic.NewForConfig(config)
	require.NoError(t, err)

	for _, obj := range template.Objs() {
		var mapping *meta.RESTMapping

		mapping, err = mapper.RESTMapping(obj.GroupVersionKind().GroupKind(), obj.GroupVersionKind().Version)
		require.NoError(t, err)

		var dr dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			// namespaced resources should specify the namespace
			dr = dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
		} else {
			// for cluster-wide resources
			dr = dyn.Resource(mapping.Resource)
		}

		var data []byte

		data, err = obj.MarshalJSON()
		require.NoError(t, err)

		t.Logf("applying %s", string(data))

		obj := obj

		_, err = dr.Create(ctx, &obj, metav1.CreateOptions{
			FieldManager: "sfyra",
		})
		if err != nil {
			if apierrors.IsAlreadyExists(err) {
				_, err = dr.Patch(ctx, obj.GetName(), types.ApplyPatchType, data, metav1.PatchOptions{
					FieldManager: "sfyra",
				})
			}
		}

		require.NoError(t, err)
	}

	t.Log("waiting for the cluster to be provisioned")

	require.NoError(t, retry.Constant(10*time.Minute, retry.WithUnits(10*time.Second), retry.WithErrorLogging(true)).Retry(func() error {
		return capi.CheckClusterReady(ctx, metalClient, clusterName)
	}))

	t.Log("verifying cluster health")

	deployedCluster, err := capi.NewCluster(ctx, metalClient, clusterName, vmSet.BridgeIP())
	require.NoError(t, err)

	require.NoError(t, deployedCluster.Health(ctx))

	return loadbalancer, deployedCluster
}

func deleteCluster(ctx context.Context, t *testing.T, metalClient client.Client, clusterName string) {
	var cluster capiv1.Cluster

	err := metalClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: clusterName}, &cluster)
	require.NoError(t, err)

	t.Logf("deleting cluster %q", clusterName)

	err = metalClient.Delete(ctx, &cluster)
	require.NoError(t, err)

	require.NoError(t, retry.Constant(3*time.Minute, retry.WithUnits(10*time.Second)).Retry(func() error {
		err = metalClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: clusterName}, &cluster)
		if err == nil {
			err = metalClient.Delete(ctx, &cluster)
			if err != nil {
				return err
			}

			return retry.ExpectedError(fmt.Errorf("cluster is not deleted yet"))
		}

		if apierrors.IsNotFound(err) {
			return nil
		}

		return err
	}))
}
