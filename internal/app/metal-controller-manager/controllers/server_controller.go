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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	goipmi "github.com/vmware/goipmi"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	metalv1alpha1 "github.com/talos-systems/sidero/internal/app/metal-controller-manager/api/v1alpha1"
)

// ServerReconciler reconciles a Server object.
type ServerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=metal.sidero.dev,resources=servers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal.sidero.dev,resources=servers/status,verbs=get;update;patch

func (r *ServerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("server", req.NamespacedName)

	s := metalv1alpha1.Server{}

	if err := r.Get(ctx, req.NamespacedName, &s); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	f := func(ready bool) (ctrl.Result, error) {
		s.Status.Ready = ready

		if err := r.Status().Update(ctx, &s); err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	if s.Spec.BMC != nil {
		conn := &goipmi.Connection{
			Hostname:  s.Spec.BMC.Endpoint,
			Username:  s.Spec.BMC.User,
			Password:  s.Spec.BMC.Pass,
			Interface: "lanplus",
		}

		client, err := goipmi.NewClient(conn)
		if err != nil {
			return f(false)
		}

		ipmiReq := &goipmi.Request{
			NetworkFunction: goipmi.NetworkFunctionChassis,
			Command:         goipmi.CommandChassisStatus,
			Data:            goipmi.ChassisStatusRequest{},
		}

		res := &goipmi.ChassisStatusResponse{}

		err = client.Send(ipmiReq, res)
		if err != nil {
			return f(false)
		}
	}

	return f(true)
}

func (r *ServerReconciler) SetupWithManager(mgr ctrl.Manager, options controller.Options) error {
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&metalv1alpha1.Server{}).
		Complete(r)
}
