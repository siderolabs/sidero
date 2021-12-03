// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	sidero "github.com/talos-systems/sidero/app/caps-controller-manager/api/v1alpha3"
)

func getMetalClient() (runtimeclient.Client, *rest.Config, error) {
	kubeconfig := ctrl.GetConfigOrDie()

	scheme := runtime.NewScheme()

	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}

	if err := sidero.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}

	client, err := runtimeclient.New(kubeconfig, runtimeclient.Options{Scheme: scheme})

	return client, kubeconfig, err
}
