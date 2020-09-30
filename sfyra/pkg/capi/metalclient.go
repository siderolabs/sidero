// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package capi

import (
	cabpt "github.com/talos-systems/cluster-api-bootstrap-provider-talos/api/v1alpha3"
	cacpt "github.com/talos-systems/cluster-api-control-plane-provider-talos/api/v1alpha3"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/cluster-api/api/v1alpha3"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	sidero "github.com/talos-systems/sidero/app/cluster-api-provider-sidero/api/v1alpha3"
	metal "github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
)

// GetMetalClient builds k8s client with schemes required to access all the CAPI/Sidero/Talos components.
func GetMetalClient(config *rest.Config) (runtimeclient.Client, error) {
	scheme := runtime.NewScheme()

	if err := v1alpha3.AddToScheme(scheme); err != nil {
		return nil, err
	}

	if err := cacpt.AddToScheme(scheme); err != nil {
		return nil, err
	}

	if err := cabpt.AddToScheme(scheme); err != nil {
		return nil, err
	}

	if err := sidero.AddToScheme(scheme); err != nil {
		return nil, err
	}

	if err := metal.AddToScheme(scheme); err != nil {
		return nil, err
	}

	return runtimeclient.New(config, runtimeclient.Options{Scheme: scheme})
}
