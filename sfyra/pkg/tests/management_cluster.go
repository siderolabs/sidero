// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tests

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/talos-systems/go-retry/retry"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"sigs.k8s.io/cluster-api/api/v1alpha3"
	capiclient "sigs.k8s.io/cluster-api/cmd/clusterctl/client"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/talos-systems/sidero/sfyra/pkg/capi"
	"github.com/talos-systems/sidero/sfyra/pkg/loadbalancer"
	"github.com/talos-systems/sidero/sfyra/pkg/talos"
	"github.com/talos-systems/sidero/sfyra/pkg/vm"
)

const (
	managementClusterName   = "management-cluster"
	managementClusterLBPort = 10000
)

// TestManagementCluster deploys the management cluster via CAPI.
//
//nolint: gocognit
func TestManagementCluster(ctx context.Context, metalClient client.Client, cluster talos.Cluster, vmSet *vm.Set, capiManager *capi.Manager) TestFunc {
	return func(t *testing.T) {
		kubeconfig, err := capiManager.GetKubeconfig(ctx)
		require.NoError(t, err)

		config, err := cluster.KubernetesClient().K8sRestConfig(ctx)
		require.NoError(t, err)

		capiClient := capiManager.GetManagerClient()

		nodeCount := int64(1)

		lb, err := loadbalancer.NewControlPlane(metalClient, vmSet.BridgeIP(), managementClusterLBPort, "default", managementClusterName, vmSet.Nodes())
		require.NoError(t, err)

		defer lb.Close()

		os.Setenv("CONTROL_PLANE_ENDPOINT", "localhost")
		os.Setenv("CONTROL_PLANE_SERVERCLASS", serverClassName)
		os.Setenv("WORKER_SERVERCLASS", serverClassName)
		// TODO: make it configurable
		os.Setenv("KUBERNETES_VERSION", "v1.19.0")

		templateOptions := capiclient.GetClusterTemplateOptions{
			Kubeconfig:               kubeconfig,
			ClusterName:              managementClusterName,
			ControlPlaneMachineCount: &nodeCount,
			WorkerMachineCount:       &nodeCount,
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

			if obj.GroupVersionKind().Kind == "MetalCluster" {
				host, portStr, _ := net.SplitHostPort(lb.GetEndpoint())
				port, _ := strconv.Atoi(portStr)

				require.NoError(t, unstructured.SetNestedMap(obj.Object, map[string]interface{}{
					"host": host,
					"port": float64(port),
				}, "spec", "controlPlaneEndpoint"))
			}

			var data []byte

			data, err = obj.MarshalJSON()
			require.NoError(t, err)

			t.Logf("applying %s", string(data))

			obj := obj

			_, err = dr.Create(ctx, &obj, metav1.CreateOptions{})
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

		require.NoError(t, retry.Constant(10*time.Minute, retry.WithUnits(10*time.Second)).Retry(func() error {
			var cluster v1alpha3.Cluster

			if err = metalClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: managementClusterName}, &cluster); err != nil {
				return retry.UnexpectedError(err)
			}

			ready := false

			for _, cond := range cluster.Status.Conditions {
				if cond.Type == v1alpha3.ReadyCondition && cond.Status == corev1.ConditionTrue {
					ready = true

					break
				}
			}

			if !ready {
				return retry.ExpectedError(fmt.Errorf("cluster is not ready"))
			}

			return nil
		}))

		t.Log("verifying cluster health")

		cluster, err := capi.NewCluster(ctx, metalClient, managementClusterName, vmSet)
		require.NoError(t, err)

		require.NoError(t, cluster.Health(ctx))
	}
}
