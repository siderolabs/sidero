// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/reference"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util/conditions"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	metalv1alpha1 "github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/internal/pkg/metal"
)

// ServerReconciler reconciles a Server object.
type ServerReconciler struct {
	client.Client
	APIReader client.Reader
	Log       logr.Logger
	Scheme    *runtime.Scheme
	Recorder  record.EventRecorder
}

// +kubebuilder:rbac:groups=metal.sidero.dev,resources=servers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal.sidero.dev,resources=servers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *ServerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("server", req.NamespacedName)

	s := metalv1alpha1.Server{}

	// Refresh the object from the API, as we rely a lot on the Status to be up to date in this controller.
	if err := r.APIReader.Get(ctx, req.NamespacedName, &s); err != nil {
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

	serverRef, err := reference.GetReference(r.Scheme, &s)
	if err != nil {
		return ctrl.Result{}, err
	}

	switch {
	case !s.Spec.Accepted:
		return f(false, false)
	case s.Status.InUse && s.Status.IsClean:
		log.Error(fmt.Errorf("server cannot be in use and clean"), "server is in an impossible state", "inUse", s.Status.InUse, "isClean", s.Status.IsClean)

		return f(false, false)
	case !s.Status.InUse && s.Status.IsClean:
		mgmtClient, err := metal.NewManagementClient(&s.Spec)
		if err != nil {
			log.Error(err, "failed to create management client")
			r.Recorder.Event(serverRef, corev1.EventTypeWarning, "Server Management", fmt.Sprintf("Failed to initialize management client: %s.", err))

			return f(false, true)
		}

		poweredOn, err := mgmtClient.IsPoweredOn()
		if err != nil {
			log.Error(err, "failed to check power state")
			r.Recorder.Event(serverRef, corev1.EventTypeWarning, "Server Management", fmt.Sprintf("Failed to determine power status: %s.", err))

			return f(false, true)
		}

		if poweredOn {
			err = mgmtClient.PowerOff()
			if err != nil {
				log.Error(err, "failed to power off")
				r.Recorder.Event(serverRef, corev1.EventTypeWarning, "Server Management", fmt.Sprintf("Failed to power off: %s.", err))

				return f(false, true)
			}

			if !mgmtClient.IsFake() {
				r.Recorder.Event(serverRef, corev1.EventTypeNormal, "Server Management", "Server powered off.")
			}
		}

		return f(true, false)
	case s.Status.InUse && !s.Status.IsClean:
		return f(true, false)
	case !s.Status.InUse && !s.Status.IsClean:
		if conditions.Has(&s, metalv1alpha1.ConditionPowerCycle) && conditions.IsFalse(&s, metalv1alpha1.ConditionPowerCycle) {
			return f(false, false)
		}

		mgmtClient, err := metal.NewManagementClient(&s.Spec)
		if err != nil {
			log.Error(err, "failed to create management client")
			r.Recorder.Event(serverRef, corev1.EventTypeWarning, "Server Management", fmt.Sprintf("Failed to initialize management client: %s.", err))

			return f(false, true)
		}

		poweredOn, err := mgmtClient.IsPoweredOn()
		if err != nil {
			log.Error(err, "failed to check power state")
			r.Recorder.Event(serverRef, corev1.EventTypeWarning, "Server Management", fmt.Sprintf("Failed to determine power status: %s.", err))

			return f(false, true)
		}

		err = mgmtClient.SetPXE()
		if err != nil {
			log.Error(err, "failed to set PXE")
			r.Recorder.Event(serverRef, corev1.EventTypeWarning, "Server Management", fmt.Sprintf("Failed to set to PXE boot once: %s.", err))

			return f(false, true)
		}

		if poweredOn {
			err = mgmtClient.PowerCycle()
			if err != nil {
				log.Error(err, "failed to power cycle")
				r.Recorder.Event(serverRef, corev1.EventTypeWarning, "Server Management", fmt.Sprintf("Failed to power cycle: %s.", err))

				return f(false, true)
			}
		} else {
			err = mgmtClient.PowerOn()
			if err != nil {
				log.Error(err, "failed to power on")
				r.Recorder.Event(serverRef, corev1.EventTypeWarning, "Server Management", fmt.Sprintf("Failed to power on: %s.", err))

				return f(false, true)
			}
		}

		if !mgmtClient.IsFake() {
			if poweredOn {
				r.Recorder.Event(serverRef, corev1.EventTypeNormal, "Server Management", "Server power cycled and set to PXE boot once.")
			} else {
				r.Recorder.Event(serverRef, corev1.EventTypeNormal, "Server Management", "Server powered on and set to PXE boot once.")
			}

			conditions.MarkFalse(&s, metalv1alpha1.ConditionPowerCycle, "InProgress", clusterv1.ConditionSeverityInfo, "Server power cycled for wiping.")
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
