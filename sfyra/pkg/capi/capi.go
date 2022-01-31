// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package capi manages CAPI installation, provides default client for CAPI CRDs.
package capi

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/talos-systems/sidero/sfyra/pkg/talos"
)

// Manager installs and controls cluster API installation.
type Manager struct {
	options Options

	cluster talos.Cluster

	kubeconfig    client.Kubeconfig
	client        client.Client
	clientset     *kubernetes.Clientset
	runtimeClient runtimeclient.Client
}

// Options for the CAPI installer.
type Options struct {
	ClusterctlConfigPath    string
	CoreProvider            string
	BootstrapProviders      []string
	InfrastructureProviders []string
	ControlPlaneProviders   []string

	PowerSimulatedExplicitFailureProb float64
	PowerSimulatedSilentFailureProb   float64
}

// NewManager creates new Manager object.
func NewManager(ctx context.Context, cluster talos.Cluster, options Options) (*Manager, error) {
	clusterAPI := &Manager{
		options: options,
		cluster: cluster,
	}

	var err error

	clusterAPI.client, err = client.New(options.ClusterctlConfigPath)
	if err != nil {
		return nil, err
	}

	clusterAPI.clientset, err = clusterAPI.cluster.KubernetesClient().K8sClient(ctx)
	if err != nil {
		return nil, err
	}

	return clusterAPI, nil
}

// GetKubeconfig returns kubeconfig in clusterctl expected format.
func (clusterAPI *Manager) GetKubeconfig(ctx context.Context) (client.Kubeconfig, error) {
	if clusterAPI.kubeconfig.Path != "" {
		return clusterAPI.kubeconfig, nil
	}

	kubeconfigBytes, err := clusterAPI.cluster.KubernetesClient().Kubeconfig(ctx)
	if err != nil {
		return client.Kubeconfig{}, err
	}

	tmpFile, err := ioutil.TempFile("", "kubeconfig")
	if err != nil {
		return client.Kubeconfig{}, err
	}

	_, err = tmpFile.Write(kubeconfigBytes)
	if err != nil {
		return client.Kubeconfig{}, err
	}

	clusterAPI.kubeconfig.Path = tmpFile.Name()
	clusterAPI.kubeconfig.Context = "admin@" + clusterAPI.cluster.Name()

	return clusterAPI.kubeconfig, nil
}

// GetManagerClient client returns instance of cluster API client.
func (clusterAPI *Manager) GetManagerClient() client.Client {
	return clusterAPI.client
}

// GetMetalClient returns k8s client stuffed with CAPI CRDs.
func (clusterAPI *Manager) GetMetalClient(ctx context.Context) (runtimeclient.Client, error) {
	if clusterAPI.runtimeClient != nil {
		return clusterAPI.runtimeClient, nil
	}

	config, err := clusterAPI.cluster.KubernetesClient().K8sRestConfig(ctx)
	if err != nil {
		return nil, err
	}

	clusterAPI.runtimeClient, err = GetMetalClient(config)

	return clusterAPI.runtimeClient, err
}

// Install the Manager components and wait for them to be ready.
func (clusterAPI *Manager) Install(ctx context.Context) error {
	kubeconfig, err := clusterAPI.GetKubeconfig(ctx)
	if err != nil {
		return err
	}

	// set template environment variables
	os.Setenv("SIDERO_CONTROLLER_MANAGER_HOST_NETWORK", "true")
	os.Setenv("SIDERO_CONTROLLER_MANAGER_DEPLOYMENT_STRATEGY", "Recreate")
	os.Setenv("SIDERO_CONTROLLER_MANAGER_API_ENDPOINT", clusterAPI.cluster.SideroComponentsIP().String())
	os.Setenv("SIDERO_CONTROLLER_MANAGER_SERVER_REBOOT_TIMEOUT", "30s") // wiping/reboot is fast in the test environment
	os.Setenv("SIDERO_CONTROLLER_MANAGER_TEST_POWER_EXPLICIT_FAILURE", fmt.Sprintf("%f", clusterAPI.options.PowerSimulatedExplicitFailureProb))
	os.Setenv("SIDERO_CONTROLLER_MANAGER_TEST_POWER_SILENT_FAILURE", fmt.Sprintf("%f", clusterAPI.options.PowerSimulatedSilentFailureProb))

	options := client.InitOptions{
		Kubeconfig:              kubeconfig,
		CoreProvider:            clusterAPI.options.CoreProvider,
		BootstrapProviders:      clusterAPI.options.BootstrapProviders,
		ControlPlaneProviders:   clusterAPI.options.ControlPlaneProviders,
		InfrastructureProviders: clusterAPI.options.InfrastructureProviders,
		TargetNamespace:         "",
		LogUsageInstructions:    false,
		WaitProviders:           true,
		WaitProviderTimeout:     5 * time.Minute,
	}

	_, err = clusterAPI.clientset.CoreV1().Namespaces().Get(ctx, "sidero-system", metav1.GetOptions{})
	if err != nil {
		_, err = clusterAPI.client.Init(options)
		return err
	}

	return nil
}
