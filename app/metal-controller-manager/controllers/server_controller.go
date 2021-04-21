// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/reference"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	infrav1 "github.com/talos-systems/sidero/app/cluster-api-provider-sidero/api/v1alpha3"
	metalv1alpha1 "github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/app/metal-controller-manager/internal/power/metal"
	"github.com/talos-systems/sidero/app/metal-controller-manager/pkg/constants"
)

const (
	serverBindingFinalizer = "storage.finalizers.server.k8s.io"
)

// ServerReconciler reconciles a Server object.
type ServerReconciler struct {
	client.Client
	Log       logr.Logger
	Scheme    *runtime.Scheme
	APIReader client.Reader
	Recorder  record.EventRecorder

	RebootTimeout time.Duration
}

// +kubebuilder:rbac:groups=metal.sidero.dev,resources=servers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal.sidero.dev,resources=servers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=serverbindings,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=serverbindings/status,verbs=get
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=metalmachines,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=metalmachines/status,verbs=get
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *ServerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("server", req.NamespacedName)

	s := metalv1alpha1.Server{}

	if err := r.APIReader.Get(ctx, req.NamespacedName, &s); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	patchHelper, err := patch.NewHelper(&s, r)
	if err != nil {
		return ctrl.Result{}, err
	}

	serverRef, err := reference.GetReference(r.Scheme, &s)
	if err != nil {
		return ctrl.Result{}, err
	}

	mgmtClient, err := metal.NewManagementClient(ctx, r.Client, &s.Spec)
	if err != nil {
		log.Error(err, "failed to create management client")
		r.Recorder.Event(serverRef, corev1.EventTypeWarning, "Server Management", fmt.Sprintf("Failed to initialize management client: %s.", err))

		return ctrl.Result{RequeueAfter: constants.DefaultRequeueAfter}, err
	}

	s.Status.Power = "off"

	poweredOn, powerErr := mgmtClient.IsPoweredOn()
	if powerErr != nil {
		s.Status.Power = "unknown"
	}

	if poweredOn {
		s.Status.Power = "on"
	}

	f := func(ready bool, result ctrl.Result) (ctrl.Result, error) {
		s.Status.Ready = ready

		if err := patchHelper.Patch(ctx, &s, patch.WithOwnedConditions{
			Conditions: []clusterv1.ConditionType{metalv1alpha1.ConditionPowerCycle, metalv1alpha1.ConditionPXEBooted},
		}); err != nil {
			return result, errors.WithStack(err)
		}

		return result, nil
	}

	allocated, serverBindingPresent, err := r.checkBinding(ctx, req)
	if err != nil {
		return ctrl.Result{}, err
	}

	if !allocated {
		if s.Status.InUse {
			// transitioning to false
			r.Recorder.Event(serverRef, corev1.EventTypeNormal, "Server Allocation", "Server marked as unallocated.")
		}

		s.Status.InUse = false

		conditions.Delete(&s, metalv1alpha1.ConditionPXEBooted)
	} else {
		s.Status.InUse = true
		s.Status.IsClean = false

		if serverBindingPresent {
			// clear any leftover ownerreferences, they were transferred by serverbinding controller
			s.OwnerReferences = []v1.OwnerReference{}
		}
	}

	hasFinalizer := controllerutil.ContainsFinalizer(&s, serverBindingFinalizer)

	if s.ObjectMeta.DeletionTimestamp.IsZero() {
		if !hasFinalizer {
			controllerutil.AddFinalizer(&s, serverBindingFinalizer)

			if err := patchHelper.Patch(ctx, &s); err != nil {
				return ctrl.Result{}, errors.WithStack(err)
			}
		}
	} else {
		// remove the finalizer from the server if it is not allocated
		if hasFinalizer && !allocated {
			controllerutil.RemoveFinalizer(&s, serverBindingFinalizer)
		}
	}

	switch {
	case !s.Spec.Accepted:
		// if server is not accepted, Sidero doesn't control server lifecycle, so we can't assume that server is (still) clean
		s.Status.IsClean = false

		return f(false, ctrl.Result{})
	case s.Status.InUse && s.Status.IsClean:
		log.Error(fmt.Errorf("server cannot be in use and clean"), "server is in an impossible state", "inUse", s.Status.InUse, "isClean", s.Status.IsClean)

		return f(false, ctrl.Result{})
	case !s.Status.InUse && s.Status.IsClean:
		if powerErr != nil {
			log.Error(powerErr, "failed to check power state")
			r.Recorder.Event(serverRef, corev1.EventTypeWarning, "Server Management", fmt.Sprintf("Failed to determine power status: %s.", powerErr))

			return f(false, ctrl.Result{RequeueAfter: constants.DefaultRequeueAfter})
		}

		if poweredOn {
			err = mgmtClient.PowerOff()
			if err != nil {
				log.Error(err, "failed to power off")
				r.Recorder.Event(serverRef, corev1.EventTypeWarning, "Server Management", fmt.Sprintf("Failed to power off: %s.", err))

				return f(false, ctrl.Result{RequeueAfter: constants.DefaultRequeueAfter})
			}

			if !mgmtClient.IsFake() {
				r.Recorder.Event(serverRef, corev1.EventTypeNormal, "Server Management", "Server powered off.")
			}
		}

		return f(true, ctrl.Result{})
	case s.Status.InUse && !s.Status.IsClean:
		if powerErr != nil {
			log.Error(powerErr, "failed to check power state")
			r.Recorder.Event(serverRef, corev1.EventTypeWarning, "Server Management", fmt.Sprintf("Failed to determine power status: %s.", powerErr))

			return f(false, ctrl.Result{RequeueAfter: constants.DefaultRequeueAfter})
		}

		if !poweredOn {
			// it's safe to set server to PXE boot even if it's already installed, as PXE server makes sure server is PXE booted only once
			err = mgmtClient.SetPXE()
			if err != nil {
				log.Error(err, "failed to set PXE")
				r.Recorder.Event(serverRef, corev1.EventTypeWarning, "Server Management", fmt.Sprintf("Failed to set to PXE boot once: %s.", err))

				return f(false, ctrl.Result{RequeueAfter: constants.DefaultRequeueAfter})
			}

			err = mgmtClient.PowerOn()
			if err != nil {
				log.Error(err, "failed to power on")
				r.Recorder.Event(serverRef, corev1.EventTypeWarning, "Server Management", fmt.Sprintf("Failed to power on: %s.", err))

				return f(false, ctrl.Result{RequeueAfter: constants.DefaultRequeueAfter})
			}

			if !mgmtClient.IsFake() {
				r.Recorder.Event(serverRef, corev1.EventTypeNormal, "Server Management", "Server powered on and set PXE boot once into the environment.")
			}
		}

		return f(true, ctrl.Result{})
	case !s.Status.InUse && !s.Status.IsClean:
		// when server is set to PXE boot to be wiped, ConditionPowerCycle is set to mark server
		// as power cycled to avoid duplicate reboot attempts from subsequent Reconciles
		//
		// we check LastTransitionTime to see if the server is in the wiping state for too long and
		// it's time to retry the IPMI sequence
		if conditions.Has(&s, metalv1alpha1.ConditionPowerCycle) &&
			conditions.IsFalse(&s, metalv1alpha1.ConditionPowerCycle) &&
			time.Since(conditions.GetLastTransitionTime(&s, metalv1alpha1.ConditionPowerCycle).Time) < r.RebootTimeout {
			// already powercycled, reboot/heartbeat timeout not elapsed, wait more
			return f(false, ctrl.Result{RequeueAfter: r.RebootTimeout / 3})
		}

		if powerErr != nil {
			log.Error(powerErr, "failed to check power state")
			r.Recorder.Event(serverRef, corev1.EventTypeWarning, "Server Management", fmt.Sprintf("Failed to determine power status: %s.", powerErr))

			return f(false, ctrl.Result{RequeueAfter: constants.DefaultRequeueAfter})
		}

		err = mgmtClient.SetPXE()
		if err != nil {
			log.Error(err, "failed to set PXE")
			r.Recorder.Event(serverRef, corev1.EventTypeWarning, "Server Management", fmt.Sprintf("Failed to set to PXE boot once: %s.", err))

			return f(false, ctrl.Result{RequeueAfter: constants.DefaultRequeueAfter})
		}

		if poweredOn {
			err = mgmtClient.PowerCycle()
			if err != nil {
				log.Error(err, "failed to power cycle")
				r.Recorder.Event(serverRef, corev1.EventTypeWarning, "Server Management", fmt.Sprintf("Failed to power cycle: %s.", err))

				return f(false, ctrl.Result{RequeueAfter: constants.DefaultRequeueAfter})
			}
		} else {
			err = mgmtClient.PowerOn()
			if err != nil {
				log.Error(err, "failed to power on")
				r.Recorder.Event(serverRef, corev1.EventTypeWarning, "Server Management", fmt.Sprintf("Failed to power on: %s.", err))

				return f(false, ctrl.Result{RequeueAfter: constants.DefaultRequeueAfter})
			}
		}

		if !mgmtClient.IsFake() {
			if poweredOn {
				r.Recorder.Event(serverRef, corev1.EventTypeNormal, "Server Management", "Server power cycled and set to PXE boot once.")
			} else {
				r.Recorder.Event(serverRef, corev1.EventTypeNormal, "Server Management", "Server powered on and set to PXE boot once.")
			}

			// remove the condition in case it was already set to make sure LastTransitionTime will be updated
			conditions.Delete(&s, metalv1alpha1.ConditionPowerCycle)
			conditions.MarkFalse(&s, metalv1alpha1.ConditionPowerCycle, "InProgress", clusterv1.ConditionSeverityInfo, "Server power cycled for wiping.")
		}

		// requeue to check for wipe timeout
		return f(false, ctrl.Result{RequeueAfter: r.RebootTimeout / 3})
	}

	return f(false, ctrl.Result{})
}

func (r *ServerReconciler) checkBinding(ctx context.Context, req ctrl.Request) (allocated, serverBindingPresent bool, err error) {
	var serverBinding infrav1.ServerBinding

	err = r.Get(ctx, req.NamespacedName, &serverBinding)
	if err == nil {
		return true, true, nil
	}

	if err != nil && !apierrors.IsNotFound(err) {
		return false, false, err
	}

	// double-check metalmachines to make sure we don't have a missing serverbinding
	var metalMachineList infrav1.MetalMachineList

	if err := r.List(ctx, &metalMachineList, client.MatchingFields(fields.Set{infrav1.MetalMachineServerRefField: req.Name})); err != nil {
		return false, false, err
	}

	for _, metalMachine := range metalMachineList.Items {
		if !metalMachine.DeletionTimestamp.IsZero() {
			continue
		}

		if metalMachine.Spec.ServerRef != nil {
			if metalMachine.Spec.ServerRef.Namespace == req.Namespace && metalMachine.Spec.ServerRef.Name == req.Name {
				return true, false, nil
			}
		}
	}

	return false, false, nil
}

func (r *ServerReconciler) SetupWithManager(mgr ctrl.Manager, options controller.Options) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &infrav1.MetalMachine{}, infrav1.MetalMachineServerRefField, func(rawObj runtime.Object) []string {
		metalMachine := rawObj.(*infrav1.MetalMachine)

		if metalMachine.Spec.ServerRef == nil {
			return nil
		}

		return []string{metalMachine.Spec.ServerRef.Name}
	}); err != nil {
		return err
	}

	mapRequests := handler.ToRequestsFunc(
		func(a handler.MapObject) []reconcile.Request {
			// servers and serverbindings always have matching names
			return []reconcile.Request{
				{
					NamespacedName: types.NamespacedName{
						Name:      a.Meta.GetName(),
						Namespace: a.Meta.GetNamespace(),
					},
				},
			}
		})

	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&metalv1alpha1.Server{}).
		Watches(
			&source.Kind{Type: &infrav1.ServerBinding{}},
			&handler.EnqueueRequestsFromMapFunc{
				ToRequests: mapRequests,
			},
		).
		Complete(r)
}
