// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package client

import (
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metalv1alpha1 "github.com/talos-systems/sidero/internal/app/metal-controller-manager/api/v1alpha1"
)

func NewClient(config *rest.Config) (c client.Client, err error) {
	scheme := runtime.NewScheme()

	if err = clientgoscheme.AddToScheme(scheme); err != nil {
		return c, err
	}

	if err = metalv1alpha1.AddToScheme(scheme); err != nil {
		return c, err
	}

	c, err = client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return c, err
	}

	return c, nil
}
