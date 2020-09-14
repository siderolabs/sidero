// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	metalv1alpha1 "github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/internal/pkg/metal"
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

	mgmtClient, err := metal.NewManagementClient(&s.Spec)
	if err != nil {
		return f(false)
	}

	_, err = mgmtClient.IsPoweredOn()
	if err != nil {
		return f(false)
	}

	return f(true)
}

func (r *ServerReconciler) SetupWithManager(mgr ctrl.Manager, options controller.Options) error {
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&metalv1alpha1.Server{}).
		Complete(r)
}
