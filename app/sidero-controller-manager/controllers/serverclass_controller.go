// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	metalv1 "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha2"
)

// ServerClassReconciler reconciles a ServerClass object.
type ServerClassReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=metal.sidero.dev,resources=serverclasses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal.sidero.dev,resources=serverclasses/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=metal.sidero.dev,resources=servers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal.sidero.dev,resources=servers/status,verbs=get;update;patch

func (r *ServerClassReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := r.Log.WithValues("serverclass", req.NamespacedName)
	l.Info("reconciling")

	//nolint:godox
	// TODO: We probably should use admission webhooks instead (or in additional) to prevent
	// unwanted edits instead of "fixing" the resource after the fact.
	if req.Name == metalv1.ServerClassAny {
		if err := ReconcileServerClassAny(ctx, r.Client); err != nil {
			return ctrl.Result{}, err
		}

		// do not return; re-reconcile it to update status
	} //nolint:wsl

	sc := metalv1.ServerClass{}

	if err := r.Get(ctx, req.NamespacedName, &sc); err != nil {
		l.Error(err, "failed fetching resource")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	patchHelper, err := patch.NewHelper(&sc, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	sl := &metalv1.ServerList{}

	if err := r.List(ctx, sl); err != nil {
		return ctrl.Result{}, fmt.Errorf("unable to get serverclass: %w", err)
	}

	results, err := metalv1.FilterServers(sl.Items,
		metalv1.AcceptedServerFilter,
		metalv1.NotCordonedServerFilter,
		sc.SelectorFilter(),
		sc.QualifiersFilter(),
	)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("unable to filter servers: %w", err)
	}

	avail := []string{}
	used := []string{}

	for _, server := range results {
		if server.Status.InUse {
			used = append(used, server.Name)
			continue
		}

		avail = append(avail, server.Name)
	}

	sc.Status.ServersAvailable = avail
	sc.Status.ServersInUse = used

	if err := patchHelper.Patch(ctx, &sc); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// ReconcileServerClassAny ensures that ServerClass "any" exist and is in desired state.
func ReconcileServerClassAny(ctx context.Context, c client.Client) error {
	key := types.NamespacedName{
		Name: metalv1.ServerClassAny,
	}

	sc := metalv1.ServerClass{}
	err := c.Get(ctx, key, &sc)

	switch {
	case apierrors.IsNotFound(err):
		sc.Name = metalv1.ServerClassAny

		return c.Create(ctx, &sc)

	case err == nil:
		patchHelper, err := patch.NewHelper(&sc, c)
		if err != nil {
			return err
		}

		sc.Spec.Qualifiers = metalv1.Qualifiers{}

		return patchHelper.Patch(ctx, &sc)

	default:
		return err
	}
}

func (r *ServerClassReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	// This mapRequests handler allows us to add a watch on server resources. Upon a server resource update,
	// we will dump all server classes and issue a reconcile against them so that they will get updated statuses
	// for available/in-use servers that match.
	mapRequests := func(a client.Object) []reconcile.Request {
		reqList := []reconcile.Request{}

		scList := &metalv1.ServerClassList{}

		if err := r.List(ctx, scList); err != nil {
			return reqList
		}

		for _, serverClass := range scList.Items {
			reqList = append(
				reqList,
				reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      serverClass.Name,
						Namespace: serverClass.Namespace,
					},
				},
			)
		}

		return reqList
	}

	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&metalv1.ServerClass{}).
		Watches(
			&source.Kind{Type: &metalv1.Server{}},
			handler.EnqueueRequestsFromMapFunc(mapRequests),
		).
		Complete(r)
}
