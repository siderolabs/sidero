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

	"github.com/siderolabs/sidero/sfyra/pkg/capi"
	"github.com/siderolabs/sidero/sfyra/pkg/talos"
)

// TestFunc is a testing function prototype.
type TestFunc func(t *testing.T)

// Options for the test.
type Options struct {
	KernelURL, InitrdURL string

	RegistryMirrors []string

	RunTestPattern string

	TalosRelease      string
	KubernetesVersion string
}

// Run all the tests.
func Run(ctx context.Context, cluster talos.Cluster, capiManager *capi.Manager, options Options) (ok bool) {
	metalClient, err := capiManager.GetMetalClient(ctx)
	if err != nil {
		log.Printf("error creating metalClient: %s", err)

		return false
	}

	testList := []testing.InternalTest{
		{
			"TestServerRegistration",
			TestServerRegistration(ctx, metalClient, cluster),
		},
		{
			"TestServerMgmtAPI",
			TestServerMgmtAPI(ctx, metalClient, cluster),
		},
		{
			"TestServerPatch",
			TestServerPatch(ctx, metalClient, options.RegistryMirrors),
		},
		{
			"TestServerValidation",
			TestServerValidation(ctx, metalClient),
		},
		{
			"TestServerAcceptance",
			TestServerAcceptance(ctx, metalClient),
		},
		{
			"TestServerCordoned",
			TestServerCordoned(ctx, metalClient),
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
			TestServerClassAny(ctx, metalClient, cluster),
		},
		{
			"TestServerClassCreate",
			TestServerClassCreate(ctx, metalClient, cluster),
		},
		{
			"TestServerClassPatch",
			TestServerClassPatch(ctx, metalClient, cluster, capiManager),
		},
		{
			"TestServerPXEBoot",
			TestServerPXEBoot(ctx, metalClient, cluster, capiManager, options.TalosRelease, options.KubernetesVersion),
		},
		{
			"TestManagementCluster",
			TestManagementCluster(ctx, metalClient, cluster, capiManager, options.TalosRelease, options.KubernetesVersion),
		},
		{
			"TestMatchServersMetalMachines",
			TestMatchServersMetalMachines(ctx, metalClient),
		},
		{
			"TestScaleWorkersUp",
			TestScaleWorkersUp(ctx, metalClient, cluster),
		},
		{
			"TestScaleWorkersDown",
			TestScaleWorkersDown(ctx, metalClient, cluster),
		},
		{
			"TestScaleControlPlaneUp",
			TestScaleControlPlaneUp(ctx, metalClient, cluster),
		},
		{
			"TestScaleControlPlaneDown",
			TestScaleControlPlaneDown(ctx, metalClient, cluster),
		},
		{
			"TestMachineDeploymentReconcile",
			TestMachineDeploymentReconcile(ctx, metalClient),
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
			TestWorkloadCluster(ctx, metalClient, cluster, capiManager, options.TalosRelease, options.KubernetesVersion),
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
