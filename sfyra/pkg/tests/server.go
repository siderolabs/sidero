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
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/sfyra/pkg/constants"
	"github.com/talos-systems/sidero/sfyra/pkg/vm"
)

// TestServerRegistration verifies that all the servers got registered.
func TestServerRegistration(ctx context.Context, metalClient client.Client, vmSet *vm.Set) TestFunc {
	return func(t *testing.T) {
		numNodes := len(vmSet.Nodes())

		var servers *v1alpha1.ServerList

		// wait for all the servers to be registered
		require.NoError(t, retry.Constant(5*time.Minute, retry.WithUnits(10*time.Second)).Retry(func() error {
			servers = &v1alpha1.ServerList{}

			if err := metalClient.List(ctx, servers); err != nil {
				return retry.UnexpectedError(err)
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
			server := v1alpha1.Server{}

			require.NoError(t, metalClient.Get(ctx, types.NamespacedName{Name: vm.UUID.String()}, &server))

			patchHelper, err := patch.NewHelper(&server, metalClient)
			require.NoError(t, err)

			server.Spec.ManagementAPI = &v1alpha1.ManagementAPI{
				Endpoint: net.JoinHostPort(bridgeIP.String(), strconv.Itoa(vm.APIPort)),
			}

			require.NoError(t, patchHelper.Patch(ctx, &server))
		}
	}
}

// TestServerPatch patches all the servers for the config.
func TestServerPatch(ctx context.Context, metalClient client.Client, talosInstaller string, registryMirrors []string) TestFunc {
	return func(t *testing.T) {
		servers := &v1alpha1.ServerList{}

		require.NoError(t, metalClient.List(ctx, servers))

		installConfig := talosconfig.InstallConfig{
			InstallDisk:       "/dev/vda",
			InstallBootloader: true,
			InstallImage:      talosInstaller,
			InstallExtraKernelArgs: []string{
				"console=ttyS0",
				"reboot=k",
				"panic=1",
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

			server.Spec.ConfigPatches = append(server.Spec.ConfigPatches, v1alpha1.ConfigPatches{
				Op:    "replace",
				Path:  "/machine/install",
				Value: apiextensions.JSON{Raw: installPatch},
			})

			if mirrorsPatch != nil {
				server.Spec.ConfigPatches = append(server.Spec.ConfigPatches, v1alpha1.ConfigPatches{
					Op:    "add",
					Path:  "/machine/registries",
					Value: apiextensions.JSON{Raw: mirrorsPatch},
				})
			}

			require.NoError(t, patchHelper.Patch(ctx, &server))
		}
	}
}

// TestServersReady waits for all the servers to be 'Ready'.
func TestServersReady(ctx context.Context, metalClient client.Client) TestFunc {
	return func(t *testing.T) {
		require.NoError(t, retry.Constant(time.Minute, retry.WithUnits(10*time.Second)).Retry(func() error {
			servers := v1alpha1.ServerList{}

			if err := metalClient.List(ctx, &servers); err != nil {
				return retry.UnexpectedError(err)
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

// createDummyServers will submit servers with dummy info that are not tied to QEMU VMs.
// A number of these, based on "count" will be created.
// These can be targeted by the spec passed in or the label "dummy-server".

//nolint: deadcode,unused
func createDummyServer(ctx context.Context, metalClient client.Client, name string, spec v1alpha1.ServerSpec) (v1alpha1.Server, error) {
	var server v1alpha1.Server

	server.APIVersion = constants.SideroAPIVersion
	server.Name = name
	server.Labels = map[string]string{"dummy-server": ""}
	server.Spec = spec

	err := metalClient.Create(ctx, &server)
	if err != nil {
		return server, err
	}

	return server, nil
}
