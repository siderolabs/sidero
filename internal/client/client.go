// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package client provides controller-runtime k8s client with all schemas.
package client

import (
	cabpt "github.com/talos-systems/cluster-api-bootstrap-provider-talos/api/v1alpha3"
	cacpt "github.com/talos-systems/cluster-api-control-plane-provider-talos/api/v1alpha3"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	caps "github.com/talos-systems/sidero/app/caps-controller-manager/api/v1alpha3"
	scm "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha1"
)

// NewClient is responsible for creating a controller-runtime k8s client with all schemas.
func NewClient(kubeconfig *string) (client.Client, error) {
	// Build rest config based on whether we've got a kubeconfig
	var (
		config *rest.Config
		err    error
	)

	if kubeconfig == nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	// Register all the schemes!
	scheme := runtime.NewScheme()

	if err = clientgoscheme.AddToScheme(scheme); err != nil {
		return nil, err
	}

	if err = capi.AddToScheme(scheme); err != nil {
		return nil, err
	}

	if err = cacpt.AddToScheme(scheme); err != nil {
		return nil, err
	}

	if err = cabpt.AddToScheme(scheme); err != nil {
		return nil, err
	}

	if err = caps.AddToScheme(scheme); err != nil {
		return nil, err
	}

	if err = scm.AddToScheme(scheme); err != nil {
		return nil, err
	}

	// Finally create client
	c, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}

	return c, nil
}
