// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	infrav1 "github.com/talos-systems/sidero/app/caps-controller-manager/api/v1alpha3"
	metalv1 "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha2"
)

// ServerBindingReconciler reconciles a ServerBinding object.
type ServerBindingReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=serverbindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=serverbindings/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=metalmachines,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=metalmachines/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=metal.sidero.dev,resources=serverclasses,verbs=get;list;watch;
// +kubebuilder:rbac:groups=metal.sidero.dev,resources=serverclasses/status,verbs=get;list;watch;
// +kubebuilder:rbac:groups=metal.sidero.dev,resources=servers,verbs=get;list;watch;
// +kubebuilder:rbac:groups=metal.sidero.dev,resources=servers/status,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *ServerBindingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, err error) {
	logger := r.Log.WithValues("serverbinding", req.NamespacedName)

	serverBinding := &infrav1.ServerBinding{}

	err = r.Get(ctx, req.NamespacedName, serverBinding)
	if apierrors.IsNotFound(err) {
		return ctrl.Result{}, nil
	}

	if err != nil {
		return ctrl.Result{}, err
	}

	// Initialize the patch helper
	patchHelper, err := patch.NewHelper(serverBinding, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Always attempt to Patch the ServerBinding object and status after each reconciliation.
	defer func() {
		if e := patchHelper.Patch(ctx, serverBinding); e != nil {
			logger.Error(e, "failed to patch metalMachine")

			if err == nil {
				err = e
			}
		}
	}()

	var server metalv1.Server

	err = r.Get(ctx, req.NamespacedName, &server)
	if err != nil {
		if apierrors.IsNotFound(err) {
			serverBinding.Status.Ready = false

			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	serverBinding.Status.Ready = true

	return ctrl.Result{}, nil
}

func (r *ServerBindingReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1.ServerBinding{}).
		Complete(r)
}
