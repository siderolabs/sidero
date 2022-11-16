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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/reference"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	infrav1 "github.com/talos-systems/sidero/app/caps-controller-manager/api/v1alpha3"
	metalv1 "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha2"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/power"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/pkg/constants"
	siderotypes "github.com/talos-systems/sidero/app/sidero-controller-manager/pkg/types"
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
	PXEMode       siderotypes.PXEMode
}

// +kubebuilder:rbac:groups=metal.sidero.dev,resources=servers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal.sidero.dev,resources=servers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=serverbindings,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=serverbindings/status,verbs=get
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=metalmachines,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=metalmachines/status,verbs=get
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

//nolint:maintidx
func (r *ServerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("server", req.NamespacedName)

	s := metalv1.Server{}

	if err := r.APIReader.Get(ctx, req.NamespacedName, &s); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	patchHelper, err := patch.NewHelper(&s, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	serverRef, err := reference.GetReference(r.Scheme, &s)
	if err != nil {
		return ctrl.Result{}, err
	}

	mgmtClient, err := power.NewManagementClient(ctx, r.Client, &s.Spec)
	if err != nil {
		log.Error(err, "failed to create management client")
		r.Recorder.Event(serverRef, corev1.EventTypeWarning, "Server Management", fmt.Sprintf("Failed to initialize management client: %s.", err))

		return ctrl.Result{RequeueAfter: constants.DefaultRequeueAfter}, err
	}

	defer mgmtClient.Close() //nolint:errcheck

	s.Status.Power = "off"

	poweredOn, powerErr := mgmtClient.IsPoweredOn()
	if powerErr != nil {
		s.Status.Power = "unknown"
	}

	if poweredOn {
		s.Status.Power = "on"
	}

	pxeMode := r.PXEMode
	if s.Spec.PXEMode != "" {
		pxeMode = s.Spec.PXEMode
	}

	f := func(ready bool, result ctrl.Result) (ctrl.Result, error) {
		s.Status.Ready = ready

		if err := patchHelper.Patch(ctx, &s, patch.WithOwnedConditions{
			Conditions: []clusterv1.ConditionType{metalv1.ConditionPowerCycle, metalv1.ConditionPXEBooted},
		}); err != nil {
			return result, errors.WithStack(err)
		}

		return result, nil
	}

	allocated, serverBinding, err := r.getServerBinding(ctx, req)
	if err != nil {
		return ctrl.Result{}, err
	}

	if !allocated {
		if s.Status.InUse {
			// transitioning to false
			r.Recorder.Event(serverRef, corev1.EventTypeNormal, "Server Allocation", "Server marked as unallocated.")
		}

		s.Status.InUse = false

		conditions.Delete(&s, metalv1.ConditionPXEBooted)
	} else {
		s.Status.InUse = true
		s.Status.IsClean = false

		if serverBinding != nil {
			// clear any leftover ownerreferences, they were transferred by serverbinding controller
			s.OwnerReferences = []v1.OwnerReference{}

			// Talos installation was successful, so mark the server as PXE booted.
			if conditions.IsTrue(serverBinding, infrav1.TalosInstalledCondition) {
				conditions.MarkTrue(serverBinding, metalv1.ConditionPXEBooted)
			}
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
			err = mgmtClient.SetPXE(pxeMode)
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

		// keep checking power state from time to time, as sometimes IPMI lies about the power state
		return f(true, ctrl.Result{RequeueAfter: constants.PowerCheckPeriod})
	case !s.Status.InUse && !s.Status.IsClean:
		// when server is set to PXE boot to be wiped, ConditionPowerCycle is set to mark server
		// as power cycled to avoid duplicate reboot attempts from subsequent Reconciles
		//
		// we check LastTransitionTime to see if the server is in the wiping state for too long and
		// it's time to retry the IPMI sequence
		if conditions.Has(&s, metalv1.ConditionPowerCycle) &&
			conditions.IsFalse(&s, metalv1.ConditionPowerCycle) &&
			time.Since(conditions.GetLastTransitionTime(&s, metalv1.ConditionPowerCycle).Time) < r.RebootTimeout {
			// already powercycled, reboot/heartbeat timeout not elapsed, wait more
			return f(false, ctrl.Result{RequeueAfter: r.RebootTimeout / 3})
		}

		if powerErr != nil {
			log.Error(powerErr, "failed to check power state")
			r.Recorder.Event(serverRef, corev1.EventTypeWarning, "Server Management", fmt.Sprintf("Failed to determine power status: %s.", powerErr))

			return f(false, ctrl.Result{RequeueAfter: constants.DefaultRequeueAfter})
		}

		err = mgmtClient.SetPXE(pxeMode)
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

			// make sure message is updated in case condition was already set to make sure LastTransitionTime will be updated
			conditions.MarkFalse(&s, metalv1.ConditionPowerCycle, "InProgress", clusterv1.ConditionSeverityInfo, fmt.Sprintf("Server power cycled for wiping at %s.", time.Now().Format(time.RFC3339)))
		}

		// requeue to check for wipe timeout
		return f(false, ctrl.Result{RequeueAfter: r.RebootTimeout / 3})
	}

	return f(false, ctrl.Result{})
}

func (r *ServerReconciler) getServerBinding(ctx context.Context, req ctrl.Request) (bool, *infrav1.ServerBinding, error) {
	var (
		serverBinding infrav1.ServerBinding
		err           error
	)

	err = r.Get(ctx, req.NamespacedName, &serverBinding)
	if err == nil {
		return true, &serverBinding, nil
	}

	if apierrors.IsNotFound(err) {
		return false, nil, nil
	}

	return false, nil, err
}

func (r *ServerReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	if err := mgr.GetFieldIndexer().IndexField(ctx, &infrav1.MetalMachine{}, infrav1.MetalMachineServerRefField, func(rawObj client.Object) []string {
		metalMachine := rawObj.(*infrav1.MetalMachine)

		if metalMachine.Spec.ServerRef == nil {
			return nil
		}

		return []string{metalMachine.Spec.ServerRef.Name}
	}); err != nil {
		return err
	}

	mapRequests := func(a client.Object) []reconcile.Request {
		// servers and serverbindings always have matching names
		return []reconcile.Request{
			{
				NamespacedName: types.NamespacedName{
					Name:      a.GetName(),
					Namespace: a.GetNamespace(),
				},
			},
		}
	}

	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&metalv1.Server{}).
		Watches(
			&source.Kind{Type: &infrav1.ServerBinding{}},
			handler.EnqueueRequestsFromMapFunc(mapRequests),
		).
		Complete(r)
}
