// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers

import (
	"context"
	"fmt"

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
	log := r.Log.WithValues("server", req.NamespacedName)

	s := metalv1alpha1.Server{}

	if err := r.Get(ctx, req.NamespacedName, &s); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	f := func(ready, requeue bool) (ctrl.Result, error) {
		s.Status.Ready = ready

		result := ctrl.Result{Requeue: requeue}

		if err := r.Status().Update(ctx, &s); err != nil {
			return result, err
		}

		return result, nil
	}

	switch {
	case s.Status.InUse && s.Status.IsClean:
		log.Error(fmt.Errorf("server cannot be in use and clean"), "server is in an impossible state", "inUse", s.Status.InUse, "isClean", s.Status.IsClean)

		return f(false, false)
	case !s.Status.InUse && s.Status.IsClean:
		mgmtClient, err := metal.NewManagementClient(&s.Spec)
		if err != nil {
			log.Error(err, "failed to create management client")

			return f(false, true)
		}

		poweredOn, err := mgmtClient.IsPoweredOn()
		if err != nil {
			log.Error(err, "failed to check power state")

			return f(false, true)
		}

		if poweredOn {
			err = mgmtClient.PowerOff()
			if err != nil {
				log.Error(err, "failed to power off")

				return f(false, true)
			}
		}

		return f(true, false)
	case s.Status.InUse && !s.Status.IsClean:
		return f(true, false)
	case !s.Status.InUse && !s.Status.IsClean:
		mgmtClient, err := metal.NewManagementClient(&s.Spec)
		if err != nil {
			log.Error(err, "failed to create management client")

			return f(false, true)
		}

		_, err = mgmtClient.IsPoweredOn()
		if err != nil {
			log.Error(err, "failed to check power state")

			return f(false, true)
		}

		err = mgmtClient.SetPXE()
		if err != nil {
			log.Error(err, "failed to set PXE")

			return f(false, true)
		}

		err = mgmtClient.PowerCycle()
		if err != nil {
			log.Error(err, "failed to power cycle")

			return f(false, true)
		}

		return f(false, false)
	}

	return f(false, false)
}

func (r *ServerReconciler) SetupWithManager(mgr ctrl.Manager, options controller.Options) error {
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&metalv1alpha1.Server{}).
		Complete(r)
}
