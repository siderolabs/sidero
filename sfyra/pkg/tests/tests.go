// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package tests provides the Sidero tests.
package tests

import (
	"context"
	"log"
	"regexp"
	"testing"

	"github.com/talos-systems/sidero/sfyra/pkg/capi"
	"github.com/talos-systems/sidero/sfyra/pkg/talos"
	"github.com/talos-systems/sidero/sfyra/pkg/vm"
)

// TestFunc is a testing function prototype.
type TestFunc func(t *testing.T)

// Options for the test.
type Options struct {
	KernelURL, InitrdURL string

	RegistryMirrors []string

	RunTestPattern string

	TalosRelease      string
	PrevTalosRelease  string
	KubernetesVersion string
}

// Run all the tests.
func Run(ctx context.Context, cluster talos.Cluster, vmSet *vm.Set, capiManager *capi.Manager, options Options) (ok bool) {
	metalClient, err := capiManager.GetMetalClient(ctx)
	if err != nil {
		log.Printf("error creating metalClient: %s", err)

		return false
	}

	testList := []testing.InternalTest{
		{
			"TestServerRegistration",
			TestServerRegistration(ctx, metalClient, vmSet),
		},
		{
			"TestServerMgmtAPI",
			TestServerMgmtAPI(ctx, metalClient, vmSet),
		},
		{
			"TestServerPatch",
			TestServerPatch(ctx, metalClient, options.RegistryMirrors),
		},
		{
			"TestServerAcceptance",
			TestServerAcceptance(ctx, metalClient, vmSet),
		},
		{
			"TestServerCordoned",
			TestServerCordoned(ctx, metalClient, vmSet),
		},
		{
			"TestServerResetOnAcceptance",
			TestServerResetOnAcceptance(ctx, metalClient),
		},
		{
			"TestServersReady",
			TestServersReady(ctx, metalClient),
		},
		{
			"TestServersDiscoveredIPs",
			TestServersDiscoveredIPs(ctx, metalClient),
		},
		{
			"TestEnvironmentDefault",
			TestEnvironmentDefault(ctx, metalClient, cluster, options.KernelURL, options.InitrdURL),
		},
		{
			"TestEnvironmentCreate",
			TestEnvironmentCreate(ctx, metalClient, cluster, options.KernelURL, options.InitrdURL),
		},
		{
			"TestServerClassAny",
			TestServerClassAny(ctx, metalClient, vmSet),
		},
		{
			"TestServerClassCreate",
			TestServerClassCreate(ctx, metalClient, vmSet),
		},
		{
			"TestServerClassPatch",
			TestServerClassPatch(ctx, metalClient, cluster, capiManager),
		},
		{
			"TestCompatibilityCluster",
			TestCompatibilityCluster(ctx, metalClient, cluster, vmSet, capiManager, options.PrevTalosRelease, "v1.22.12"),
		},
		{
			"TestServerPXEBoot",
			TestServerPXEBoot(ctx, metalClient, cluster, vmSet, capiManager, options.TalosRelease, options.KubernetesVersion),
		},
		{
			"TestManagementCluster",
			TestManagementCluster(ctx, metalClient, cluster, vmSet, capiManager, options.TalosRelease, options.KubernetesVersion),
		},
		{
			"TestMatchServersMetalMachines",
			TestMatchServersMetalMachines(ctx, metalClient),
		},
		{
			"TestScaleWorkersUp",
			TestScaleWorkersUp(ctx, metalClient, vmSet),
		},
		{
			"TestScaleWorkersDown",
			TestScaleWorkersDown(ctx, metalClient, vmSet),
		},
		{
			"TestScaleControlPlaneUp",
			TestScaleControlPlaneUp(ctx, metalClient, vmSet),
		},
		{
			"TestScaleControlPlaneDown",
			TestScaleControlPlaneDown(ctx, metalClient, vmSet),
		},
		{
			"TestMachineDeploymentReconcile",
			TestMachineDeploymentReconcile(ctx, metalClient),
		},
		{
			"TestServerBindingReconcile",
			TestServerBindingReconcile(ctx, metalClient),
		},
		{
			"TestMetalMachineServerRefReconcile",
			TestMetalMachineServerRefReconcile(ctx, metalClient),
		},
		{
			"TestServerReset",
			TestServerReset(ctx, metalClient),
		},
		{
			"TestWorkloadCluster",
			TestWorkloadCluster(ctx, metalClient, cluster, vmSet, capiManager, options.TalosRelease, options.KubernetesVersion),
		},
	}

	testsToRun := []testing.InternalTest{}

	var re *regexp.Regexp

	if options.RunTestPattern != "" {
		if re, err = regexp.Compile(options.RunTestPattern); err != nil {
			log.Printf("run test pattern parse error: %s", err)

			return false
		}
	}

	for _, test := range testList {
		if re == nil || re.MatchString(test.Name) {
			testsToRun = append(testsToRun, test)
		}
	}

	return testing.MainStart(matchStringOnly(func(pat, str string) (bool, error) { return true, nil }), testsToRun, nil, nil, nil).Run() == 0
}
