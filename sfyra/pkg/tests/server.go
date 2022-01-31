// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/talos-systems/go-retry/retry"
	talosconfig "github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "github.com/talos-systems/sidero/app/caps-controller-manager/api/v1alpha3"
	metalv1 "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha2"
	"github.com/talos-systems/sidero/sfyra/pkg/capi"
	"github.com/talos-systems/sidero/sfyra/pkg/constants"
	"github.com/talos-systems/sidero/sfyra/pkg/talos"
	"github.com/talos-systems/sidero/sfyra/pkg/vm"
)

// TestServerRegistration verifies that all the servers got registered.
func TestServerRegistration(ctx context.Context, metalClient client.Client, vmSet *vm.Set) TestFunc {
	return func(t *testing.T) {
		numNodes := len(vmSet.Nodes())

		var servers *metalv1.ServerList

		// wait for all the servers to be registered
		require.NoError(t, retry.Constant(5*time.Minute, retry.WithUnits(10*time.Second)).Retry(func() error {
			servers = &metalv1.ServerList{}

			if err := metalClient.List(ctx, servers); err != nil {
				return err
			}

			if len(servers.Items) != numNodes {
				return retry.ExpectedError(fmt.Errorf("%d != %d", len(servers.Items), numNodes))
			}

			return nil
		}))

		assert.Len(t, servers.Items, numNodes)

		nodes := vmSet.Nodes()
		expectedUUIDs := make([]string, len(nodes))

		for i := range nodes {
			expectedUUIDs[i] = nodes[i].UUID.String()
		}

		actualUUIDs := make([]string, len(servers.Items))

		for i := range servers.Items {
			actualUUIDs[i] = servers.Items[i].Name
		}

		sort.Strings(expectedUUIDs)
		sort.Strings(actualUUIDs)

		assert.Equal(t, expectedUUIDs, actualUUIDs)
	}
}

func configPatchToJSON(t *testing.T, o interface{}) []byte {
	patchYaml, err := yaml.Marshal(o)
	require.NoError(t, err)

	var obj map[string]interface{}

	require.NoError(t, yaml.Unmarshal(patchYaml, &obj))

	patchJSON, err := json.Marshal(obj)
	require.NoError(t, err)

	return patchJSON
}

// TestServerMgmtAPI patches all the servers for the management API.
func TestServerMgmtAPI(ctx context.Context, metalClient client.Client, vmSet *vm.Set) TestFunc {
	return func(t *testing.T) {
		bridgeIP := vmSet.BridgeIP()

		for _, vm := range vmSet.Nodes() {
			server := metalv1.Server{}

			require.NoError(t, metalClient.Get(ctx, types.NamespacedName{Name: vm.UUID.String()}, &server))

			patchHelper, err := patch.NewHelper(&server, metalClient)
			require.NoError(t, err)

			server.Spec.ManagementAPI = &metalv1.ManagementAPI{
				Endpoint: net.JoinHostPort(bridgeIP.String(), strconv.Itoa(vm.APIPort)),
			}

			require.NoError(t, patchHelper.Patch(ctx, &server))
		}
	}
}

// TestServerPatch patches all the servers for the config.
func TestServerPatch(ctx context.Context, metalClient client.Client, registryMirrors []string) TestFunc {
	return func(t *testing.T) {
		servers := &metalv1.ServerList{}

		require.NoError(t, metalClient.List(ctx, servers))

		installConfig := talosconfig.InstallConfig{
			InstallDisk:       "/dev/vda",
			InstallBootloader: true,
			InstallExtraKernelArgs: []string{
				"console=ttyS0",
				"reboot=k",
				"panic=1",
				"talos.shutdown=halt",
			},
		}
		installPatch := configPatchToJSON(t, &installConfig)

		var mirrorsPatch []byte

		if len(registryMirrors) > 0 {
			var registriesConfig talosconfig.RegistriesConfig

			registriesConfig.RegistryMirrors = make(map[string]*talosconfig.RegistryMirrorConfig)

			for _, mirror := range registryMirrors {
				parts := strings.SplitN(mirror, "=", 2)
				require.Len(t, parts, 2)

				registriesConfig.RegistryMirrors[parts[0]] = &talosconfig.RegistryMirrorConfig{
					MirrorEndpoints: []string{parts[1]},
				}
			}

			mirrorsPatch = configPatchToJSON(t, &registriesConfig)
		}

		for _, server := range servers.Items {
			if len(server.Spec.ConfigPatches) > 0 {
				continue
			}

			server := server

			patchHelper, err := patch.NewHelper(&server, metalClient)
			require.NoError(t, err)

			server.Spec.ConfigPatches = append(server.Spec.ConfigPatches, metalv1.ConfigPatches{
				Op:    "replace",
				Path:  "/machine/install",
				Value: apiextensions.JSON{Raw: installPatch},
			})

			if mirrorsPatch != nil {
				server.Spec.ConfigPatches = append(server.Spec.ConfigPatches, metalv1.ConfigPatches{
					Op:    "add",
					Path:  "/machine/registries",
					Value: apiextensions.JSON{Raw: mirrorsPatch},
				})
			}

			require.NoError(t, patchHelper.Patch(ctx, &server))
		}
	}
}

// TestServerAcceptance makes sure the accepted bool works.
func TestServerAcceptance(ctx context.Context, metalClient client.Client, vmSet *vm.Set) TestFunc {
	return func(t *testing.T) {
		const numDummies = 3

		// create dummy servers to test with
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
		}

		for i := 0; i < numDummies; i++ {
			serverName := fmt.Sprintf("dummyserver-%s", strconv.Itoa(i))
			_, err := createDummyServer(ctx, metalClient, serverName, dummySpec)
			require.NoError(t, err)
		}

		dummyServers := &metalv1.ServerList{}

		labelSelector, err := labels.Parse("dummy-server=")
		require.NoError(t, err)

		// wait for all the servers to be registered
		require.NoError(t, retry.Constant(1*time.Minute, retry.WithUnits(10*time.Second)).Retry(func() error {
			err = metalClient.List(ctx, dummyServers, client.MatchingLabelsSelector{Selector: labelSelector})
			if err != nil {
				return err
			}

			if len(dummyServers.Items) != numDummies {
				return retry.ExpectedError(fmt.Errorf("%d != %d", len(dummyServers.Items), numDummies))
			}

			return nil
		}))

		// verify servers originally registered as non-accepted
		acceptedServers := []metalv1.Server{}

		for _, server := range dummyServers.Items {
			if server.Spec.Accepted {
				acceptedServers = append(acceptedServers, server)
			}
		}

		assert.Len(t, acceptedServers, 0)

		// patch all servers as accepted
		for _, server := range dummyServers.Items {
			server := server

			patchHelper, err := patch.NewHelper(&server, metalClient)
			require.NoError(t, err)

			server.Spec.Accepted = true
			require.NoError(t, patchHelper.Patch(ctx, &server))
		}

		// verify all servers are now accepted
		require.NoError(t, metalClient.List(ctx, dummyServers, client.MatchingLabelsSelector{Selector: labelSelector}))

		acceptedServers = []metalv1.Server{}

		for _, server := range dummyServers.Items {
			if server.Spec.Accepted {
				acceptedServers = append(acceptedServers, server)
			}
		}

		assert.Len(t, acceptedServers, numDummies)

		// clean up dummies
		for _, server := range dummyServers.Items {
			server := server
			require.NoError(t, metalClient.Delete(ctx, &server))
		}
	}
}

// TestServerCordoned makes sure the cordoned bool works.
func TestServerCordoned(ctx context.Context, metalClient client.Client, vmSet *vm.Set) TestFunc {
	return func(t *testing.T) {
		const numDummies = 3

		// create dummy servers to test with
		dummySpec := metalv1.ServerSpec{
			Hardware: &metalv1.HardwareInformation{
				Compute: &metalv1.ComputeInformation{
					Processors: []*metalv1.Processor{
						{
							Manufacturer: "DummyManufacturer",
						},
					},
				},
			},
		}

		for i := 0; i < numDummies; i++ {
			serverName := fmt.Sprintf("dummyserver-%s", strconv.Itoa(i))
			_, err := createDummyServer(ctx, metalClient, serverName, dummySpec)
			require.NoError(t, err)
		}

		dummyServers := &metalv1.ServerList{}

		labelSelector, err := labels.Parse("dummy-server=")
		require.NoError(t, err)
		err = metalClient.List(ctx, dummyServers, client.MatchingLabelsSelector{Selector: labelSelector})
		require.NoError(t, err)

		// clean up dummies
		defer func(client client.Client) {
			for _, server := range dummyServers.Items {
				server := server
				client.Delete(ctx, &server)
			}
		}(metalClient)

		// patch all servers as accepted
		for _, server := range dummyServers.Items {
			server := server

			patchHelper, err := patch.NewHelper(&server, metalClient)
			require.NoError(t, err)

			server.Spec.Accepted = true
			require.NoError(t, patchHelper.Patch(ctx, &server))
		}

		// verify that all servers shows up as available in `any` serverclass
		require.NoError(t, retry.Constant(30*time.Second, retry.WithUnits(5*time.Second)).Retry(func() error {
			var serverClass metalv1.ServerClass
			err := metalClient.Get(ctx, types.NamespacedName{Name: metalv1.ServerClassAny}, &serverClass)
			if err != nil {
				return err
			}

			availableServers := getAvailableServersFromServerClass(serverClass, dummyServers)
			if len(availableServers) == numDummies {
				return nil
			}

			return retry.ExpectedError(fmt.Errorf("%d != %d", len(availableServers), numDummies))
		}))

		// // cordon a single server and marked as paused
		serverName := dummyServers.Items[0].Name

		var server metalv1.Server

		require.NoError(t, metalClient.Get(ctx, types.NamespacedName{Name: serverName}, &server))
		patchHelper, err := patch.NewHelper(&server, metalClient)
		require.NoError(t, err)

		server.Spec.Cordoned = true

		require.NoError(t, patchHelper.Patch(ctx, &server))

		require.NoError(t, retry.Constant(30*time.Second, retry.WithUnits(5*time.Second)).Retry(func() error {
			var serverClass metalv1.ServerClass
			err := metalClient.Get(ctx, types.NamespacedName{Name: metalv1.ServerClassAny}, &serverClass)
			if err != nil {
				return err
			}

			availableServers := getAvailableServersFromServerClass(serverClass, dummyServers)
			if len(availableServers) == numDummies-1 {
				return nil
			}

			return retry.ExpectedError(fmt.Errorf("%d != %d", len(availableServers), numDummies-1))
		}))

		// patch the server and marked as not cordoned
		var pausedServer metalv1.Server

		require.NoError(t, metalClient.Get(ctx, types.NamespacedName{Name: serverName}, &pausedServer))
		patchHelperPausedServer, err := patch.NewHelper(&pausedServer, metalClient)
		require.NoError(t, err)

		pausedServer.Spec.Cordoned = false

		require.NoError(t, patchHelperPausedServer.Patch(ctx, &pausedServer))

		require.NoError(t, retry.Constant(30*time.Second, retry.WithUnits(5*time.Second)).Retry(func() error {
			var serverClass metalv1.ServerClass
			err := metalClient.Get(ctx, types.NamespacedName{Name: metalv1.ServerClassAny}, &serverClass)
			if err != nil {
				return err
			}

			availableServers := getAvailableServersFromServerClass(serverClass, dummyServers)
			if len(availableServers) == numDummies {
				return nil
			}

			return retry.ExpectedError(fmt.Errorf("%d != %d", len(availableServers), numDummies))
		}))
	}
}

// TestServerResetOnAcceptance tests that servers are reset when accepted.
func TestServerResetOnAcceptance(ctx context.Context, metalClient client.Client) TestFunc {
	return func(t *testing.T) {
		serverList := &metalv1.ServerList{}

		err := metalClient.List(ctx, serverList)
		require.NoError(t, err)

		servers := []metalv1.Server{}

		for _, server := range serverList.Items {
			server := server

			if !server.Spec.Accepted {
				patchHelper, err := patch.NewHelper(&server, metalClient)
				require.NoError(t, err)

				server.Spec.Accepted = true
				require.NoError(t, patchHelper.Patch(ctx, &server))

				servers = append(servers, server)
			}
		}

		t.Logf("Found %d dirty servers", len(servers))

		if len(servers) == 0 {
			return
		}

		require.NoError(t, retry.Constant(5*time.Minute, retry.WithUnits(10*time.Second)).Retry(func() error {
			for _, server := range servers {
				var s metalv1.Server

				if err := metalClient.Get(ctx, types.NamespacedName{Name: server.Name, Namespace: server.Namespace}, &s); err != nil {
					return err
				}

				if !s.Status.IsClean {
					return retry.ExpectedError(fmt.Errorf("server %q is not clean", s.Name))
				}

				t.Logf("Server %q is clean", s.Name)
			}

			return nil
		}))
	}
}

// TestServersReady waits for all the servers to be 'Ready'.
func TestServersReady(ctx context.Context, metalClient client.Client) TestFunc {
	return func(t *testing.T) {
		require.NoError(t, retry.Constant(time.Minute, retry.WithUnits(10*time.Second)).Retry(func() error {
			servers := metalv1.ServerList{}

			if err := metalClient.List(ctx, &servers); err != nil {
				return err
			}

			for _, server := range servers.Items {
				if !server.Status.Ready {
					return retry.ExpectedError(fmt.Errorf("server %q is not ready", server.Name))
				}
			}

			return nil
		}))
	}
}

// TestServersDiscoveredIPs waits for all the servers to have an IP address.
func TestServersDiscoveredIPs(ctx context.Context, metalClient client.Client) TestFunc {
	return func(t *testing.T) {
		require.NoError(t, retry.Constant(time.Minute, retry.WithUnits(10*time.Second)).Retry(func() error {
			servers := metalv1.ServerList{}

			if err := metalClient.List(ctx, &servers); err != nil {
				return err
			}

			for _, server := range servers.Items {
				found := false

				for _, address := range server.Status.Addresses {
					if address.Type == corev1.NodeInternalIP {
						found = true

						break
					}
				}

				if !found {
					return retry.ExpectedError(fmt.Errorf("server %q doesn't have an internal IP address", server.Name))
				}
			}

			return nil
		}))
	}
}

const (
	pxeTestClusterName   = "pxe-test-cluster"
	pxeTestClusterLBPort = 10002
)

// TestServerPXEBoot verifies that PXE boot is retried when the server gets incorrect configuration.
func TestServerPXEBoot(ctx context.Context, metalClient client.Client, cluster talos.Cluster, vmSet *vm.Set, capiManager *capi.Manager, talosRelease, kubernetesVersion string) TestFunc {
	return func(t *testing.T) {
		pxeTestServerClass := "pxe-test-server"

		classSpec := metalv1.ServerClassSpec{
			Qualifiers: metalv1.Qualifiers{
				Hardware: []metalv1.HardwareInformation{
					{
						System: &metalv1.SystemInformation{
							Manufacturer: "QEMU",
						},
					},
				},
			},
			EnvironmentRef: &v1.ObjectReference{
				Name: environmentName,
			},
			ConfigPatches: []metalv1.ConfigPatches{
				{
					Op:    "add",
					Path:  "/fake",
					Value: apiextensions.JSON{Raw: []byte("\":|\"")},
				},
			},
		}

		serverClass, err := createServerClass(ctx, metalClient, pxeTestServerClass, classSpec)
		require.NoError(t, err)

		loadbalancer := createCluster(ctx, t, metalClient, cluster, vmSet, capiManager, pxeTestClusterName, pxeTestServerClass, pxeTestClusterLBPort, 1, 0, talosRelease, kubernetesVersion)

		retry.Constant(time.Minute, retry.WithUnits(10*time.Second)).Retry(func() error {
			var machines infrav1.MetalMachineList

			labelSelector, err := labels.Parse(fmt.Sprintf("cluster.x-k8s.io/cluster-name=%s", pxeTestClusterName))
			require.NoError(t, err)

			err = metalClient.List(ctx, &machines, client.MatchingLabelsSelector{Selector: labelSelector})
			if err != nil {
				return fmt.Errorf("failed refetching dummy server: %w", err)
			}

			if len(machines.Items) == 0 {
				return retry.ExpectedErrorf("no metal machines detected yet")
			}

			if !conditions.IsFalse(&machines.Items[0], infrav1.TalosConfigLoadedCondition) || !conditions.IsFalse(&machines.Items[0], infrav1.TalosConfigValidatedCondition) {
				return retry.ExpectedErrorf("the machine doesn't have any config failure conditions yet")
			}

			return nil
		})

		patchHelper, err := patch.NewHelper(&serverClass, metalClient)
		require.NoError(t, err)

		serverClass.Spec.ConfigPatches = nil

		err = patchHelper.Patch(ctx, &serverClass)
		require.NoError(t, err)

		waitForClusterReady(ctx, t, metalClient, vmSet, pxeTestClusterName)

		deleteCluster(ctx, t, metalClient, pxeTestClusterName)
		loadbalancer.Close() //nolint:errcheck
	}
}

// createDummyServers will submit servers with dummy info that are not tied to QEMU VMs.
// These can be targeted by the spec passed in or the label "dummy-server".
// Dummy servers are patched after creation to ensure they're marked as clean.
func createDummyServer(ctx context.Context, metalClient client.Client, name string, spec metalv1.ServerSpec) (metalv1.Server, error) {
	var server metalv1.Server

	server.APIVersion = constants.SideroAPIVersion
	server.Name = name
	server.Labels = map[string]string{"dummy-server": ""}
	server.Spec = spec

	err := metalClient.Create(ctx, &server)
	if err != nil {
		return server, fmt.Errorf("failed creating dummy server: %w", err)
	}

	return server, retry.Constant(time.Minute, retry.WithUnits(10*time.Second)).Retry(func() error {
		// refetch dummy server to make sure we're synced up before patching
		server = metalv1.Server{}

		err = metalClient.Get(ctx, types.NamespacedName{Name: name}, &server)
		if err != nil {
			return fmt.Errorf("failed refetching dummy server: %w", err)
		}

		patchHelper, err := patch.NewHelper(&server, metalClient)
		if err != nil {
			return fmt.Errorf("failed creating patch helper for dummy server: %w", err)
		}

		server.Status.InUse = false
		server.Status.IsClean = true
		server.Status.Ready = true

		err = patchHelper.Patch(ctx, &server)
		if err != nil {
			return retry.ExpectedError(fmt.Errorf("failed patching dummy server: %w", err))
		}

		return nil
	})
}

// getAvailableServersFromServerClass returns a list of servers that are available as part of a serverclass.
func getAvailableServersFromServerClass(serverClass metalv1.ServerClass, serverList *metalv1.ServerList) []string {
	var foundServers []string

	for _, server := range serverList.Items {
		for _, serverName := range serverClass.Status.ServersAvailable {
			if server.Name == serverName {
				foundServers = append(foundServers, serverName)
			}
		}
	}

	return foundServers
}
