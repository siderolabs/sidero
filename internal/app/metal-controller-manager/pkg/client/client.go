/*
Copyright 2020 Talos Systems, Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
