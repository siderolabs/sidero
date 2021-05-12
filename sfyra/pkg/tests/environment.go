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
	"github.com/talos-systems/go-procfs/procfs"
	"github.com/talos-systems/go-retry/retry"
	"github.com/talos-systems/talos/pkg/machinery/kernel"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/sfyra/pkg/constants"
	"github.com/talos-systems/sidero/sfyra/pkg/talos"
)

const environmentName = "default"

// TestEnvironmentDefault verifies default environment creation.
func TestEnvironmentDefault(ctx context.Context, metalClient client.Client, cluster talos.Cluster, kernelURL, initrdURL string) TestFunc {
	return func(t *testing.T) {
		var environment v1alpha1.Environment

		if err := metalClient.Get(ctx, types.NamespacedName{Name: environmentName}, &environment); err != nil {
			if !apierrors.IsNotFound(err) {
				require.NoError(t, err)
			}

			cmdline := procfs.NewCmdline("")
			cmdline.SetAll(kernel.DefaultArgs)

			cmdline.Append("console", "ttyS0")
			cmdline.Append("reboot", "k")
			cmdline.Append("panic", "1")
			cmdline.Append("talos.platform", "metal")
			cmdline.Append("talos.shutdown", "halt")
			cmdline.Append("talos.config", fmt.Sprintf("http://%s:9091/configdata?uuid=", cluster.SideroComponentsIP()))
			cmdline.Append("initrd", "initramfs.xz")

			environment.APIVersion = constants.SideroAPIVersion
			environment.Name = environmentName
			environment.Spec.Kernel.URL = kernelURL
			environment.Spec.Kernel.SHA512 = "" // TODO: add a test
			environment.Spec.Kernel.Args = cmdline.Strings()
			environment.Spec.Initrd.URL = initrdURL
			environment.Spec.Initrd.SHA512 = "" // TODO: add a test

			require.NoError(t, metalClient.Create(ctx, &environment))
		}

		// wait for the environment to report ready
		require.NoError(t, retry.Constant(5*time.Minute, retry.WithUnits(10*time.Second)).Retry(func() error {
			if err := metalClient.Get(ctx, types.NamespacedName{Name: environmentName}, &environment); err != nil {
				return retry.UnexpectedError(err)
			}

			assetURLs := map[string]struct{}{
				kernelURL: {},
				initrdURL: {},
			}

			for _, cond := range environment.Status.Conditions {
				if cond.Status == "True" && cond.Type == "Ready" {
					delete(assetURLs, cond.URL)
				}
			}

			if len(assetURLs) > 0 {
				return retry.ExpectedError(fmt.Errorf("some assets are not ready: %v", assetURLs))
			}

			return nil
		}))
	}
}
