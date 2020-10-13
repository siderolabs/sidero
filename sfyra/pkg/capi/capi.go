// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package capi manages CAPI installation, provides default client for CAPI CRDs.
package capi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
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
	BootstrapProviders      []string
	InfrastructureProviders []string
	ControlPlaneProviders   []string
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

	options := client.InitOptions{
		Kubeconfig:              kubeconfig,
		CoreProvider:            "",
		BootstrapProviders:      clusterAPI.options.BootstrapProviders,
		ControlPlaneProviders:   clusterAPI.options.ControlPlaneProviders,
		InfrastructureProviders: clusterAPI.options.InfrastructureProviders,
		TargetNamespace:         "",
		WatchingNamespace:       "",
		LogUsageInstructions:    false,
	}

	_, err = clusterAPI.clientset.CoreV1().Namespaces().Get(ctx, "sidero-system", metav1.GetOptions{})
	if err != nil {
		_, err = clusterAPI.client.Init(options)
		if err != nil {
			return err
		}
	}

	return clusterAPI.patch(ctx)
}

func (clusterAPI *Manager) patch(ctx context.Context) error {
	const (
		sideroNamespace         = "sidero-system"
		sideroMetadataServer    = "sidero-metadata-server"
		sideroControllerManager = "sidero-controller-manager"
	)

	// sidero-metadata-server
	deployment, err := clusterAPI.clientset.AppsV1().Deployments(sideroNamespace).Get(ctx, sideroMetadataServer, metav1.GetOptions{})
	if err != nil {
		return err
	}

	oldDeployment, err := json.Marshal(deployment)
	if err != nil {
		return err
	}

	argsPatched := false

	for _, arg := range deployment.Spec.Template.Spec.Containers[0].Args {
		if arg == "--port=9091" {
			argsPatched = true
		}
	}

	if !argsPatched {
		deployment.Spec.Template.Spec.Containers[0].Args = append(deployment.Spec.Template.Spec.Containers[0].Args, "--port=9091")
	}

	deployment.Spec.Template.Spec.Containers[0].Ports = []corev1.ContainerPort{
		{
			ContainerPort: 9091,
			HostPort:      9091,
			Name:          "http",
			Protocol:      corev1.ProtocolTCP,
		},
	}
	deployment.Spec.Template.Spec.HostNetwork = true
	deployment.Spec.Strategy.RollingUpdate = nil
	deployment.Spec.Strategy.Type = appsv1.RecreateDeploymentStrategyType

	newDeployment, err := json.Marshal(deployment)
	if err != nil {
		return err
	}

	patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldDeployment, newDeployment, appsv1.Deployment{})
	if err != nil {
		return fmt.Errorf("failed to create two way merge patch: %w", err)
	}

	_, err = clusterAPI.clientset.AppsV1().Deployments(sideroNamespace).Patch(ctx, deployment.Name, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{
		FieldManager: "sfyra",
	})
	if err != nil {
		return err
	}

	// sidero-controller-manager
	deployment, err = clusterAPI.clientset.AppsV1().Deployments(sideroNamespace).Get(ctx, sideroControllerManager, metav1.GetOptions{})
	if err != nil {
		return err
	}

	oldDeployment, err = json.Marshal(deployment)
	if err != nil {
		return err
	}

	apiPatch := false

	for _, arg := range deployment.Spec.Template.Spec.Containers[1].Args {
		if strings.HasPrefix(arg, "--api-endpoint") {
			apiPatch = true
		}
	}

	if !apiPatch {
		deployment.Spec.Template.Spec.Containers[1].Args = append(
			deployment.Spec.Template.Spec.Containers[1].Args,
			fmt.Sprintf("--api-endpoint=%s", clusterAPI.cluster.SideroComponentsIP()),
		)
	}

	deployment.Spec.Template.Spec.HostNetwork = true
	deployment.Spec.Strategy.RollingUpdate = nil
	deployment.Spec.Strategy.Type = appsv1.RecreateDeploymentStrategyType

	newDeployment, err = json.Marshal(deployment)
	if err != nil {
		return err
	}

	patchBytes, err = strategicpatch.CreateTwoWayMergePatch(oldDeployment, newDeployment, appsv1.Deployment{})
	if err != nil {
		return fmt.Errorf("failed to create two way merge patch: %w", err)
	}

	_, err = clusterAPI.clientset.AppsV1().Deployments(sideroNamespace).Patch(ctx, deployment.Name, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{
		FieldManager: "sfyra",
	})
	if err != nil {
		return err
	}

	return nil
}
