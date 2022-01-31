// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tests

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/stretchr/testify/require"
	"github.com/talos-systems/go-procfs/procfs"
	"github.com/talos-systems/go-retry/retry"
	"github.com/talos-systems/talos/pkg/machinery/kernel"

	metalv1 "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha2"
	"github.com/talos-systems/sidero/sfyra/pkg/capi"
	"github.com/talos-systems/sidero/sfyra/pkg/constants"
	"github.com/talos-systems/sidero/sfyra/pkg/talos"
	"github.com/talos-systems/sidero/sfyra/pkg/vm"
)

const (
	compatibilityClusterName   = "compatibility-cluster"
	compatibilityClusterLBPort = 10003
)

// TestCompatibilityCluster deploys the compatibility cluster via CAPI.
func TestCompatibilityCluster(ctx context.Context, metalClient client.Client, cluster talos.Cluster, vmSet *vm.Set, capiManager *capi.Manager, talosRelease, kubernetesVersion string) TestFunc {
	return func(t *testing.T) {
		if talosRelease == "" {
			t.Skip("--prev-talos-release is not set, skipped compatibility check")
		}

		var environment metalv1.Environment

		envName := fmt.Sprintf("talos-%s", strings.ReplaceAll(talosRelease, ".", "-"))

		if err := metalClient.Get(ctx, types.NamespacedName{Name: envName}, &environment); err != nil {
			if !apierrors.IsNotFound(err) {
				require.NoError(t, err)
			}

			cmdline := procfs.NewCmdline("")
			cmdline.SetAll(kernel.DefaultArgs)

			cmdline.Append("console", "ttyS0")
			cmdline.Append("talos.platform", "metal")

			environment.APIVersion = constants.SideroAPIVersion
			environment.Name = envName
			environment.Spec.Kernel.URL = fmt.Sprintf("https://github.com/talos-systems/talos/releases/download/%s/vmlinuz-amd64", talosRelease)
			environment.Spec.Kernel.SHA512 = ""
			environment.Spec.Kernel.Args = cmdline.Strings()
			environment.Spec.Initrd.URL = fmt.Sprintf("https://github.com/talos-systems/talos/releases/download/%s/initramfs-amd64.xz", talosRelease)
			environment.Spec.Initrd.SHA512 = ""

			require.NoError(t, metalClient.Create(ctx, &environment))
		}

		// wait for the environment to report ready
		require.NoError(t, retry.Constant(5*time.Minute, retry.WithUnits(10*time.Second)).Retry(func() error {
			if err := metalClient.Get(ctx, types.NamespacedName{Name: envName}, &environment); err != nil {
				return err
			}

			if !environment.IsReady() {
				return retry.ExpectedErrorf("some assets are not ready")
			}

			return nil
		}))

		serverClassName := envName
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
				Name: envName,
			},
		}

		_, err := createServerClass(ctx, metalClient, serverClassName, classSpec)
		require.NoError(t, err)

		ex, err := os.Executable()
		require.NoError(t, err)

		exPath := filepath.Dir(ex)

		loadbalancer := deployCluster(ctx, t, metalClient, cluster, vmSet, capiManager, compatibilityClusterName, serverClassName, compatibilityClusterLBPort, 1, 0, talosRelease, kubernetesVersion,
			withConfigURL(fmt.Sprintf("file://%s/../templates/cluster-template-talos-%s.yaml", exPath, talosRelease)),
		)

		deleteCluster(ctx, t, metalClient, compatibilityClusterName)
		loadbalancer.Close() //nolint:errcheck
	}
}
