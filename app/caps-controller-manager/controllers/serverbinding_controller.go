// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

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
		return r.reconcileTransition(ctx, logger, req)
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
	if err := mgr.GetFieldIndexer().IndexField(ctx, &infrav1.MetalMachine{}, infrav1.MetalMachineServerRefField, func(rawObj client.Object) []string {
		metalMachine := rawObj.(*infrav1.MetalMachine)

		if metalMachine.Spec.ServerRef == nil {
			return nil
		}

		return []string{metalMachine.Spec.ServerRef.Name}
	}); err != nil {
		return err
	}

	// This mapRequests handler allows us to transition to the new scheme with ServerBinding.
	mapRequests := func(a client.Object) []reconcile.Request {
		metalMachine := &infrav1.MetalMachine{}

		if err := r.Get(context.Background(), types.NamespacedName{Namespace: a.GetNamespace(), Name: a.GetName()}, metalMachine); err != nil {
			return nil
		}

		if metalMachine.Spec.ServerRef == nil {
			return nil
		}

		return []reconcile.Request{
			{
				NamespacedName: types.NamespacedName{
					Name:      metalMachine.Spec.ServerRef.Name,
					Namespace: metalMachine.Spec.ServerRef.Namespace,
				},
			},
		}
	}

	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1.ServerBinding{}).
		Watches(
			&source.Kind{Type: &infrav1.MetalMachine{}},
			handler.EnqueueRequestsFromMapFunc(mapRequests),
		).
		Complete(r)
}

func (r *ServerBindingReconciler) reconcileTransition(ctx context.Context, logger logr.Logger, req ctrl.Request) (_ ctrl.Result, err error) {
	logger.Info("reconciling missing serverbinding")

	var metalMachineList infrav1.MetalMachineList

	if err := r.List(ctx, &metalMachineList, client.MatchingFields(fields.Set{infrav1.MetalMachineServerRefField: req.Name})); err != nil {
		return ctrl.Result{}, err
	}

	var serverBinding infrav1.ServerBinding

	serverBinding.Namespace = req.Namespace
	serverBinding.Name = req.Name
	serverBinding.Labels = map[string]string{}

	found := false

	for _, metalMachine := range metalMachineList.Items {
		if !metalMachine.DeletionTimestamp.IsZero() {
			continue
		}

		if metalMachine.Spec.ServerRef != nil {
			if metalMachine.Spec.ServerRef.Name == serverBinding.Name && metalMachine.Spec.ServerRef.Namespace == serverBinding.Namespace {
				found = true

				serverBinding.Spec.MetalMachineRef = corev1.ObjectReference{
					Kind:      metalMachine.Kind,
					UID:       metalMachine.UID,
					Namespace: metalMachine.Namespace,
					Name:      metalMachine.Name,
				}

				if metalMachine.Spec.ServerClassRef != nil {
					serverBinding.Spec.ServerClassRef = metalMachine.Spec.ServerClassRef.DeepCopy()
				}

				for label, value := range metalMachine.Labels {
					serverBinding.Labels[label] = value
				}

				break
			}
		}
	}

	if !found {
		logger.Info("no matching metalmachine found")

		return ctrl.Result{}, nil
	}

	var server metalv1.Server

	if err = r.Get(ctx, req.NamespacedName, &server); err != nil {
		if apierrors.IsNotFound(err) {
			// broken link?
			logger.Info("server not found", "name", req.Name)

			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	for _, ownerRef := range server.OwnerReferences {
		if ownerRef.Kind == "ServerClass" {
			serverBinding.Spec.ServerClassRef = &corev1.ObjectReference{
				Kind: ownerRef.Kind,
				Name: ownerRef.Name,
			}
		}
	}

	logger.Info("creating missing server binding", "metalmachine", serverBinding.Spec.MetalMachineRef.Name)

	err = r.Create(ctx, &serverBinding)
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			err = nil
		}
	}

	return ctrl.Result{}, err
}
