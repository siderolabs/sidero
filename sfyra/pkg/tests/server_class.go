// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tests

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/talos-systems/go-retry/retry"
	"github.com/talos-systems/talos/pkg/machinery/config/configloader"
	talosconfig "github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1"
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

	"github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/sfyra/pkg/capi"
	"github.com/talos-systems/sidero/sfyra/pkg/constants"
	"github.com/talos-systems/sidero/sfyra/pkg/talos"
	"github.com/talos-systems/sidero/sfyra/pkg/vm"
)

const (
	defaultServerClassName  = "default"
	workloadServerClassName = "workload"
)

func TestServerClassAny(ctx context.Context, metalClient client.Client, vmSet *vm.Set) TestFunc {
	return func(t *testing.T) {
		var serverClass v1alpha1.ServerClass
		err := metalClient.Get(ctx, types.NamespacedName{Name: v1alpha1.ServerClassAny}, &serverClass)
		require.NoError(t, err)
		assert.Empty(t, serverClass.Spec.Qualifiers)

		numNodes := len(vmSet.Nodes())
		assert.Len(t, append(serverClass.Status.ServersAvailable, serverClass.Status.ServersInUse...), numNodes)
	}
}

// TestServerClassDefault verifies server class creation.
func TestServerClassDefault(ctx context.Context, metalClient client.Client, vmSet *vm.Set) TestFunc {
	return func(t *testing.T) {
		classSpec := v1alpha1.ServerClassSpec{
			Qualifiers: v1alpha1.Qualifiers{
				CPU: []v1alpha1.CPUInformation{
					{
						Manufacturer: "QEMU",
					},
				},
			},
		}

		serverClass, err := createServerClass(ctx, metalClient, defaultServerClassName, classSpec)
		require.NoError(t, err)

		numNodes := len(vmSet.Nodes())

		// wait for the server class to gather all nodes (all nodes should match)
		require.NoError(t, retry.Constant(2*time.Minute, retry.WithUnits(10*time.Second)).Retry(func() error {
			if err := metalClient.Get(ctx, types.NamespacedName{Name: defaultServerClassName}, &serverClass); err != nil {
				return retry.UnexpectedError(err)
			}

			if len(serverClass.Status.ServersAvailable)+len(serverClass.Status.ServersInUse) != numNodes {
				return retry.ExpectedError(fmt.Errorf("%d + %d != %d", len(serverClass.Status.ServersAvailable), len(serverClass.Status.ServersInUse), numNodes))
			}

			return nil
		}))

		assert.Len(t, append(serverClass.Status.ServersAvailable, serverClass.Status.ServersInUse...), numNodes)

		nodes := vmSet.Nodes()
		expectedUUIDs := make([]string, len(nodes))

		for i := range nodes {
			expectedUUIDs[i] = nodes[i].UUID.String()
		}

		actualUUIDs := append(serverClass.Status.ServersAvailable, serverClass.Status.ServersInUse...)

		sort.Strings(expectedUUIDs)
		sort.Strings(actualUUIDs)

		assert.Equal(t, expectedUUIDs, actualUUIDs)

		_, err = createServerClass(ctx, metalClient, workloadServerClassName, classSpec)
		require.NoError(t, err)
	}
}

// TestServerClassPatch verifies config patches work at the server level.
func TestServerClassPatch(ctx context.Context, metalClient client.Client, cluster talos.Cluster, capiManager *capi.Manager) TestFunc {
	return func(t *testing.T) {
		// Create dummy serverclass + a server
		dummySpec := v1alpha1.ServerSpec{
			CPU: &v1alpha1.CPUInformation{
				Manufacturer: "DummyCPU",
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

		classSpec := v1alpha1.ServerClassSpec{
			ConfigPatches: []v1alpha1.ConfigPatches{
				{
					Op:    "replace",
					Path:  "/machine/install",
					Value: apiextensions.JSON{Raw: installPatch},
				},
			},
			Qualifiers: v1alpha1.Qualifiers{
				CPU: []v1alpha1.CPUInformation{
					{
						Manufacturer: "DummyCPU",
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

		os.Setenv("CONTROL_PLANE_ENDPOINT", "localhost")
		os.Setenv("CONTROL_PLANE_PORT", "11111")
		os.Setenv("CONTROL_PLANE_SERVERCLASS", "dummyservers")
		os.Setenv("WORKER_SERVERCLASS", "dummyservers")
		// TODO: make it configurable
		os.Setenv("KUBERNETES_VERSION", "v1.20.4")
		os.Setenv("TALOS_VERSION", "v0.9")

		templateOptions := capiclient.GetClusterTemplateOptions{
			Kubeconfig:               kubeconfig,
			ClusterName:              "serverclass-config-patch-test",
			ControlPlaneMachineCount: &nodeCountCP,
			WorkerMachineCount:       &nodeCountWorker,
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

			obj := obj

			_, err = dr.Create(ctx, &obj, metav1.CreateOptions{})
			require.NoError(t, err)

			if obj.GroupVersionKind().Kind == "Cluster" {
				//nolint:errcheck
				defer dr.Delete(ctx, obj.GetName(), metav1.DeleteOptions{})
			}
		}

		// Wait for metadata server to return a 200 for that UUID, then verify it has patch.
		metadataEndpoint := fmt.Sprintf("http://%s:9091/configdata?uuid=dummyserver-0", cluster.SideroComponentsIP())

		var metadataBytes []byte

		require.NoError(t, retry.Constant(5*time.Minute, retry.WithUnits(10*time.Second)).Retry(func() error {
			req, _ := http.NewRequestWithContext(ctx, "GET", metadataEndpoint, nil)
			client := &http.Client{}
			response, err := client.Do(req)
			if err != nil {
				return retry.UnexpectedError(err)
			}

			if response.StatusCode != http.StatusOK {
				return retry.ExpectedError(fmt.Errorf("metadata not yet present: %d", response.StatusCode))
			}

			defer response.Body.Close()
			metadataBytes, err = ioutil.ReadAll(response.Body)
			if err != nil {
				return retry.UnexpectedError(err)
			}
			return nil
		}))

		configProvider, err := configloader.NewFromBytes(metadataBytes)
		require.NoError(t, err)

		// try deleting dummy server before deleting the cluster
		// it shouldn't break the cluster
		require.NoError(t, metalClient.Delete(ctx, &dummyServer))

		time.Sleep(time.Second * 10)

		response := &v1alpha1.Server{}

		require.NoError(t, metalClient.Get(ctx, types.NamespacedName{Name: dummyServer.Name}, response))
		require.Greater(t, len(response.Finalizers), 0)

		switch configProvider.Version() {
		case "v1alpha1":
			config, ok := configProvider.(*talosconfig.Config)
			if !ok {
				t.Error("unable to case config")
			}

			require.Len(t, config.MachineConfig.MachineInstall.InstallExtraKernelArgs, 1)

			if config.MachineConfig.MachineInstall.InstallExtraKernelArgs[0] != "woot" {
				t.Error("unable to validate server class patch was applied to kernel args")
			}

		default:
			t.Error("unknown config type")
		}
	}
}

func createServerClass(ctx context.Context, metalClient client.Client, name string, spec v1alpha1.ServerClassSpec) (v1alpha1.ServerClass, error) {
	var retClass v1alpha1.ServerClass

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
