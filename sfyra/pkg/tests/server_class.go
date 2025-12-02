// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tests

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"
	"testing"
	"time"

	"github.com/siderolabs/go-retry/retry"
	"github.com/siderolabs/talos/pkg/machinery/config/configloader"
	talosconfig "github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	capiclient "sigs.k8s.io/cluster-api/cmd/clusterctl/client"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metalv1 "github.com/siderolabs/sidero/app/sidero-controller-manager/api/v1alpha2"
	"github.com/siderolabs/sidero/sfyra/pkg/capi"
	"github.com/siderolabs/sidero/sfyra/pkg/constants"
	"github.com/siderolabs/sidero/sfyra/pkg/talos"
)

const (
	defaultServerClassName = "default"
	serverClassName        = "sfyra"
)

// TestServerClassAny verifies server class "any".
func TestServerClassAny(ctx context.Context, metalClient client.Client, capiCluster talos.Cluster) TestFunc {
	return func(t *testing.T) {
		numNodes := len(capiCluster.Nodes())

		var serverClass metalv1.ServerClass

		err := metalClient.Get(ctx, types.NamespacedName{Name: metalv1.ServerClassAny}, &serverClass)
		require.NoError(t, err)
		assert.Empty(t, serverClass.Spec.Qualifiers)
		assert.Len(t, append(serverClass.Status.ServersAvailable, serverClass.Status.ServersInUse...), numNodes)

		// delete server class to see it being recreated
		err = metalClient.Delete(ctx, &serverClass)
		require.NoError(t, err)

		serverClass = metalv1.ServerClass{}
		err = retry.Constant(10 * time.Second).Retry(func() error {
			if err := metalClient.Get(ctx, types.NamespacedName{Name: metalv1.ServerClassAny}, &serverClass); err != nil {
				if apierrors.IsNotFound(err) {
					return retry.ExpectedError(err)
				}

				return err
			}

			if len(serverClass.Status.ServersAvailable)+len(serverClass.Status.ServersInUse) != numNodes {
				return retry.ExpectedErrorf("%d + %d != %d", len(serverClass.Status.ServersAvailable), len(serverClass.Status.ServersInUse), numNodes)
			}

			return nil
		})
		require.NoError(t, err)
		assert.Empty(t, serverClass.Spec.Qualifiers)
		assert.Len(t, append(serverClass.Status.ServersAvailable, serverClass.Status.ServersInUse...), numNodes)
	}
}

// TestServerClassCreate verifies server class creation.
func TestServerClassCreate(ctx context.Context, metalClient client.Client, capiCluster talos.Cluster) TestFunc {
	return func(t *testing.T) {
		classSpec := metalv1.ServerClassSpec{
			Qualifiers: metalv1.Qualifiers{
				Hardware: []metalv1.HardwareInformation{
					{
						Compute: &metalv1.ComputeInformation{
							Processors: []*metalv1.Processor{
								{
									Manufacturer: "QEMU",
								},
							},
						},
					},
				},
			},
			EnvironmentRef: &v1.ObjectReference{
				Name: environmentName,
			},
		}

		serverClass, err := createServerClass(ctx, metalClient, defaultServerClassName, classSpec)
		require.NoError(t, err)

		numNodes := len(capiCluster.Nodes())

		// wait for the server class to gather all nodes (all nodes should match)
		require.NoError(t, retry.Constant(2*time.Minute, retry.WithUnits(10*time.Second)).Retry(func() error {
			if err := metalClient.Get(ctx, types.NamespacedName{Name: defaultServerClassName}, &serverClass); err != nil {
				return err
			}

			if len(serverClass.Status.ServersAvailable)+len(serverClass.Status.ServersInUse) != numNodes {
				return retry.ExpectedErrorf("%d + %d != %d", len(serverClass.Status.ServersAvailable), len(serverClass.Status.ServersInUse), numNodes)
			}

			return nil
		}))

		assert.Len(t, append(serverClass.Status.ServersAvailable, serverClass.Status.ServersInUse...), numNodes)

		nodes := capiCluster.Nodes()
		expectedUUIDs := make([]string, len(nodes))

		for i := range nodes {
			expectedUUIDs[i] = nodes[i].UUID.String()
		}

		actualUUIDs := append(append([]string(nil), serverClass.Status.ServersAvailable...), serverClass.Status.ServersInUse...)

		sort.Strings(expectedUUIDs)
		sort.Strings(actualUUIDs)

		assert.Equal(t, expectedUUIDs, actualUUIDs)

		_, err = createServerClass(ctx, metalClient, serverClassName, classSpec)
		require.NoError(t, err)
	}
}

// TestServerClassPatch verifies config patches work at the server level.
func TestServerClassPatch(ctx context.Context, metalClient client.Client, cluster talos.Cluster, capiManager *capi.Manager) TestFunc {
	return func(t *testing.T) {
		// Create dummy serverclass + a server
		dummySpec := metalv1.ServerSpec{
			Hardware: &metalv1.HardwareInformation{
				Compute: &metalv1.ComputeInformation{
					Processors: []*metalv1.Processor{
						{
							Manufacturer: "DummyCPU",
						},
					},
				},
			},
			Accepted: true,
		}

		dummyServer, err := createDummyServer(ctx, metalClient, "dummyserver-0", dummySpec)
		require.NoError(t, err)

		//nolint:errcheck
		defer metalClient.Delete(ctx, &dummyServer)

		installConfig := talosconfig.InstallConfig{
			InstallExtraKernelArgs: []string{
				"woot",
			},
		}

		installPatch := configPatchToJSON(t, &installConfig)

		classSpec := metalv1.ServerClassSpec{
			ConfigPatches: []metalv1.ConfigPatches{
				{
					Op:    "replace",
					Path:  "/machine/install",
					Value: apiextensions.JSON{Raw: installPatch},
				},
			},
			Qualifiers: metalv1.Qualifiers{
				Hardware: []metalv1.HardwareInformation{
					{
						Compute: &metalv1.ComputeInformation{
							Processors: []*metalv1.Processor{
								{
									Manufacturer: "DummyCPU",
								},
							},
						},
					},
				},
			},
		}

		dummyClass, err := createServerClass(ctx, metalClient, "dummyservers", classSpec)
		require.NoError(t, err)

		//nolint:errcheck
		defer metalClient.Delete(ctx, &dummyClass)

		// Create "cluster" using serverclass above
		kubeconfig, err := capiManager.GetKubeconfig(ctx)
		require.NoError(t, err)

		config, err := cluster.KubernetesClient().K8sRestConfig(ctx)
		require.NoError(t, err)

		capiClient := capiManager.GetManagerClient()

		nodeCountCP := int64(1)
		nodeCountWorker := int64(0)

		t.Setenv("CONTROL_PLANE_ENDPOINT", "localhost")
		t.Setenv("CONTROL_PLANE_PORT", "11111")
		t.Setenv("CONTROL_PLANE_SERVERCLASS", "dummyservers")
		t.Setenv("WORKER_SERVERCLASS", "dummyservers")
		t.Setenv("KUBERNETES_VERSION", "v1.20.4") // dummy cluster, actual value doesn't matter
		t.Setenv("TALOS_VERSION", "v0.9")         // dummy cluster, actual value doesn't matter

		templateOptions := capiclient.GetClusterTemplateOptions{
			Kubeconfig:               kubeconfig,
			ClusterName:              "serverclass-config-patch-test",
			ControlPlaneMachineCount: &nodeCountCP,
			WorkerMachineCount:       &nodeCountWorker,
		}

		template, err := capiClient.GetClusterTemplate(ctx, templateOptions)
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

			obj := obj

			_, err = dr.Create(ctx, &obj, metav1.CreateOptions{})
			require.NoError(t, err)

			if obj.GroupVersionKind().Kind == "Cluster" {
				//nolint:errcheck
				defer dr.Delete(ctx, obj.GetName(), metav1.DeleteOptions{})
			}
		}

		// Wait for metadata server to return a 200 for that UUID, then verify it has patch.
		metadataEndpoint := fmt.Sprintf("http://%s/configdata?uuid=dummyserver-0", net.JoinHostPort(cluster.SideroComponentsIP().String(), "8081"))

		var metadataBytes []byte

		require.NoError(t, retry.Constant(5*time.Minute, retry.WithUnits(10*time.Second)).Retry(func() error {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, metadataEndpoint, nil)
			if err != nil {
				return err
			}

			client := &http.Client{}

			response, err := client.Do(req)
			if err != nil {
				return err
			}

			if response.StatusCode != http.StatusOK {
				return retry.ExpectedErrorf("metadata not yet present: %d", response.StatusCode)
			}

			defer response.Body.Close()

			metadataBytes, err = io.ReadAll(response.Body)
			if err != nil {
				return err
			}

			return nil
		}))

		configProvider, err := configloader.NewFromBytes(metadataBytes)
		require.NoError(t, err)

		// try deleting dummy server before deleting the cluster
		// it shouldn't break the cluster
		require.NoError(t, metalClient.Delete(ctx, &dummyServer))

		time.Sleep(time.Second * 10)

		response := &metalv1.Server{}

		require.NoError(t, metalClient.Get(ctx, types.NamespacedName{Name: dummyServer.Name}, response))
		require.Greater(t, len(response.Finalizers), 0)

		configV1Alpha1 := configProvider.RawV1Alpha1()

		if config == nil {
			t.Error("unable to cast config")
		}

		require.Len(t, configV1Alpha1.MachineConfig.MachineInstall.InstallExtraKernelArgs, 1)

		if configV1Alpha1.MachineConfig.MachineInstall.InstallExtraKernelArgs[0] != "woot" {
			t.Error("unable to validate server class patch was applied to kernel args")
		}
	}
}

func createServerClass(ctx context.Context, metalClient client.Client, name string, spec metalv1.ServerClassSpec) (metalv1.ServerClass, error) {
	var retClass metalv1.ServerClass

	if err := metalClient.Get(ctx, types.NamespacedName{Name: name}, &retClass); err != nil {
		if !apierrors.IsNotFound(err) {
			return retClass, nil
		}

		retClass.APIVersion = constants.SideroAPIVersion
		retClass.Name = name
		retClass.Spec = spec

		err = metalClient.Create(ctx, &retClass)
		if err != nil {
			return retClass, err
		}
	}

	return retClass, nil
}
