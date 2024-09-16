// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/siderolabs/go-pointer"
	talosclient "github.com/siderolabs/talos/pkg/machinery/client"
	clientconfig "github.com/siderolabs/talos/pkg/machinery/client/config"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/reference"
	capiv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/remote"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "github.com/siderolabs/sidero/app/caps-controller-manager/api/v1alpha3"
	"github.com/siderolabs/sidero/app/caps-controller-manager/pkg/constants"
	metalv1 "github.com/siderolabs/sidero/app/sidero-controller-manager/api/v1alpha2"
)

var ErrNoServersInServerClass = errors.New("no servers available in serverclass")

// MetalMachineReconciler reconciles a MetalMachine object.
type MetalMachineReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Tracker  *remote.ClusterCacheTracker
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=metalmachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=metalmachines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machines/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=metal.sidero.dev,resources=serverclasses,verbs=get;list;watch;
// +kubebuilder:rbac:groups=metal.sidero.dev,resources=serverclasses/status,verbs=get;list;watch;
// +kubebuilder:rbac:groups=metal.sidero.dev,resources=servers,verbs=get;list;watch;
// +kubebuilder:rbac:groups=metal.sidero.dev,resources=servers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *MetalMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, err error) {
	logger := r.Log.WithValues("metalmachine", req.NamespacedName)

	// Fetch the metalMachine instance.
	metalMachine := &infrav1.MetalMachine{}

	err = r.Get(ctx, req.NamespacedName, metalMachine)
	if apierrors.IsNotFound(err) {
		return ctrl.Result{}, nil
	}

	if err != nil {
		return ctrl.Result{}, err
	}

	machine, err := util.GetOwnerMachine(ctx, r.Client, metalMachine.ObjectMeta)
	if err != nil {
		r.Log.Error(err, "Failed to get machine")

		return ctrl.Result{RequeueAfter: constants.DefaultRequeueAfter}, nil
	}

	if machine == nil {
		logger.Info("No ownerref for metalmachine")

		return ctrl.Result{RequeueAfter: constants.DefaultRequeueAfter}, nil
	}

	logger = logger.WithName(fmt.Sprintf("machine=%s", machine.Name))

	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, machine.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("no cluster label or cluster does not exist")
	}

	logger = logger.WithName(fmt.Sprintf("cluster=%s", cluster.Name))

	if !cluster.Status.InfrastructureReady {
		logger.Error(err, "Cluster infrastructure is not ready", "cluster", cluster.Name)

		return ctrl.Result{RequeueAfter: constants.DefaultRequeueAfter}, nil
	}

	if machine.Spec.Bootstrap.DataSecretName == nil {
		logger.Info("Bootstrap secret is not available yet")

		return ctrl.Result{RequeueAfter: constants.DefaultRequeueAfter}, nil
	}

	// Initialize the patch helper
	patchHelper, err := patch.NewHelper(metalMachine, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Always attempt to Patch the MetalMachine object and status after each reconciliation.
	defer func() {
		if e := patchHelper.Patch(ctx, metalMachine); e != nil {
			logger.Error(e, "failed to patch metalMachine")

			if err == nil {
				err = e
			}
		}
	}()

	// Handle deleted machines
	if !metalMachine.ObjectMeta.DeletionTimestamp.IsZero() {
		logger.Info("deleting metalmachine")

		return r.reconcileDelete(ctx, metalMachine)
	}

	controllerutil.AddFinalizer(metalMachine, infrav1.MachineFinalizer)

	// TODO (smira):
	// This is really weird that .Spec.ServerRef is used both to set the manual link to the server by the user
	// and a way to store current binding when server is picked up automatically
	//
	// This should be refactored to keep _current_ binding ref in the .Status, while
	// .spec should be reserved for manually chosen server only.
	// This opens a question what to do if the `.Status` is lost after pivoting, how to reconcile it correctly?

	// if server binding is missing, need to pick up a server
	if metalMachine.Spec.ServerRef == nil {
		if metalMachine.Spec.ServerClassRef == nil {
			return ctrl.Result{}, fmt.Errorf("either a server or serverclass ref must be supplied")
		}

		serverResource, err := r.fetchServerFromClass(ctx, logger, metalMachine.Spec.ServerClassRef, metalMachine)
		if err != nil {
			if errors.Is(err, ErrNoServersInServerClass) {
				return ctrl.Result{RequeueAfter: constants.DefaultRequeueAfter}, nil
			}

			return ctrl.Result{}, err
		}

		metalMachine.Spec.ServerRef = &corev1.ObjectReference{
			Kind: serverResource.Kind,
			Name: serverResource.Name,
		}
	} else {
		// If server ref is set, it could have been set by "us", or by the user.
		// In any case, check if we need to create server binding
		var serverBinding infrav1.ServerBinding

		err = r.Get(ctx, types.NamespacedName{Namespace: metalMachine.Spec.ServerRef.Namespace, Name: metalMachine.Spec.ServerRef.Name}, &serverBinding)
		if err != nil {
			if apierrors.IsNotFound(err) {
				var serverObj metalv1.Server

				namespacedName := types.NamespacedName{
					Namespace: "",
					Name:      metalMachine.Spec.ServerRef.Name,
				}

				if err := r.Get(ctx, namespacedName, &serverObj); err != nil {
					return ctrl.Result{}, err
				}

				serverRef, err := reference.GetReference(r.Scheme, &serverObj)
				if err != nil {
					return ctrl.Result{}, err
				}

				if err = r.createServerBinding(ctx, nil, &serverObj, metalMachine); err != nil {
					return ctrl.Result{}, err
				}

				r.Recorder.Event(serverRef, corev1.EventTypeNormal, "Server Assignment", fmt.Sprintf("Server as assigned via serverRef for metal machine %q.", metalMachine.Name))
			}

			return ctrl.Result{}, err
		}
	}

	// Set the providerID, as its required in upstream capi for machine lifecycle
	metalMachine.Spec.ProviderID = pointer.To(fmt.Sprintf("%s://%s", constants.ProviderID, metalMachine.Spec.ServerRef.Name))

	// Copy over statuses from ServerBinding to MetalMachine
	if metalMachine.Spec.ServerRef != nil {
		var serverBinding infrav1.ServerBinding

		err = r.Get(ctx, types.NamespacedName{Namespace: metalMachine.Spec.ServerRef.Namespace, Name: metalMachine.Spec.ServerRef.Name}, &serverBinding)
		if err != nil {
			if apierrors.IsNotFound(err) {
				return ctrl.Result{RequeueAfter: constants.DefaultRequeueAfter}, nil
			}

			return ctrl.Result{}, err
		}

		addresses := make([]capiv1.MachineAddress, 0, len(serverBinding.Spec.Addresses))
		for _, addr := range serverBinding.Spec.Addresses {
			addresses = append(addresses, capiv1.MachineAddress{
				Type:    capiv1.MachineInternalIP,
				Address: addr,
			})
		}

		if serverBinding.Spec.Hostname != "" {
			addresses = append(addresses, capiv1.MachineAddress{
				Type:    capiv1.MachineHostName,
				Address: serverBinding.Spec.Hostname,
			})
		}

		metalMachine.Status.Addresses = addresses
		metalMachine.Status.Ready = true

		// copy conditions from the server binding
		for _, condition := range serverBinding.GetConditions() {
			conditions.Set(metalMachine, &condition)
		}
	}

	err = r.patchProviderID(ctx, cluster, metalMachine)
	if err != nil {
		logger.Info("Failed to set provider ID", "error", err)

		conditions.MarkFalse(metalMachine, infrav1.ProviderSetCondition, infrav1.ProviderUpdateFailedReason, capiv1.ConditionSeverityWarning, err.Error())

		return ctrl.Result{RequeueAfter: constants.DefaultRequeueAfter}, nil
	}

	conditions.MarkTrue(metalMachine, infrav1.ProviderSetCondition)

	return ctrl.Result{}, nil
}

func (r *MetalMachineReconciler) reconcileDelete(ctx context.Context, metalMachine *infrav1.MetalMachine) (ctrl.Result, error) {
	if metalMachine.Spec.ServerRef != nil {
		var serverBinding infrav1.ServerBinding

		err := r.Get(ctx, types.NamespacedName{Namespace: metalMachine.Spec.ServerRef.Namespace, Name: metalMachine.Spec.ServerRef.Name}, &serverBinding)
		if err == nil {
			if err = r.ResetServer(ctx, metalMachine, &serverBinding); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, r.Delete(ctx, &serverBinding)
		}

		if !apierrors.IsNotFound(err) {
			return ctrl.Result{}, err
		}
	}

	metalMachine.Spec.ServerRef = nil

	controllerutil.RemoveFinalizer(metalMachine, infrav1.MachineFinalizer)

	return ctrl.Result{}, nil
}

func (r *MetalMachineReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	if err := mgr.GetFieldIndexer().IndexField(ctx, &infrav1.ServerBinding{}, infrav1.ServerBindingMetalMachineRefField, func(rawObj client.Object) []string {
		serverBinding := rawObj.(*infrav1.ServerBinding)

		return []string{serverBinding.Spec.MetalMachineRef.Name}
	}); err != nil {
		return err
	}

	mapRequests := func(ctx context.Context, a client.Object) []reconcile.Request {
		serverBinding := &infrav1.ServerBinding{}

		if err := r.Get(ctx, types.NamespacedName{Namespace: a.GetNamespace(), Name: a.GetName()}, serverBinding); err != nil {
			return nil
		}

		return []reconcile.Request{
			{
				NamespacedName: types.NamespacedName{
					Name:      serverBinding.Spec.MetalMachineRef.Name,
					Namespace: serverBinding.Spec.MetalMachineRef.Namespace,
				},
			},
		}
	}

	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1.MetalMachine{}).
		Watches(
			&infrav1.ServerBinding{},
			handler.EnqueueRequestsFromMapFunc(mapRequests),
		).
		Complete(r)
}

func (r *MetalMachineReconciler) fetchServerFromClass(ctx context.Context, logger logr.Logger, classRef *corev1.ObjectReference, metalMachine *infrav1.MetalMachine) (*metalv1.Server, error) {
	// First, check if there is already existing serverBinding for this metalmachine
	var serverBindingList infrav1.ServerBindingList

	if err := r.List(ctx, &serverBindingList, client.MatchingFields(fields.Set{infrav1.ServerBindingMetalMachineRefField: metalMachine.Name})); err != nil {
		return nil, err
	}

	for _, serverBinding := range serverBindingList.Items {
		if serverBinding.Spec.MetalMachineRef.Namespace == metalMachine.Namespace && serverBinding.Spec.MetalMachineRef.Name == metalMachine.Name {
			// found existing serverBinding for this metalMachine
			var server metalv1.Server

			if err := r.Get(ctx, types.NamespacedName{Namespace: serverBinding.Namespace, Name: serverBinding.Name}, &server); err != nil {
				return nil, err
			}

			logger.Info("reconciled missing server ref", "metalmachine", metalMachine.Name, "server", server.Name)

			return &server, nil
		}
	}

	// Grab server class and validate that we have nodes available
	serverClassResource, err := r.fetchServerClass(ctx, classRef)
	if err != nil {
		return nil, err
	}

	if len(serverClassResource.Status.ServersAvailable) == 0 {
		return nil, ErrNoServersInServerClass
	}

	serverClassRef, err := reference.GetReference(r.Scheme, serverClassResource)
	if err != nil {
		return nil, err
	}

	// Fetch server from available list
	// NB: we added this loop to double check that an available server isn't "in use" because
	//     we saw raciness between server selection and it being removed from the ServersAvailable list.
	for _, availServer := range serverClassResource.Status.ServersAvailable {
		serverObj := &metalv1.Server{}

		namespacedName := types.NamespacedName{
			Namespace: "",
			Name:      availServer,
		}

		if err := r.Get(ctx, namespacedName, serverObj); err != nil {
			return nil, err
		}

		serverRef, err := reference.GetReference(r.Scheme, serverObj)
		if err != nil {
			return nil, err
		}

		if serverObj.Status.InUse {
			continue
		}

		if !serverObj.Status.IsClean {
			continue
		}

		if err := r.createServerBinding(ctx, serverClassRef, serverObj, metalMachine); err != nil {
			// the server we picked was updated by another metalmachine before we finished.
			// move on to the next one.
			if apierrors.IsAlreadyExists(err) {
				continue
			}

			return nil, err
		}

		r.Recorder.Event(serverRef, corev1.EventTypeNormal, "Server Allocation", fmt.Sprintf("Server is allocated via serverclass %q for metal machine %q.", serverClassResource.Name, metalMachine.Name))

		logger.Info("allocated new server", "metalmachine", metalMachine.Name, "server", serverObj.Name, "serverclass", serverClassResource.Name)

		return serverObj, nil
	}

	return nil, ErrNoServersInServerClass
}

func (r *MetalMachineReconciler) patchProviderID(ctx context.Context, cluster *capiv1.Cluster, metalMachine *infrav1.MetalMachine) error {
	workloadClient, err := r.Tracker.GetClient(ctx, client.ObjectKeyFromObject(cluster))
	if err != nil {
		return err
	}

	var nodes corev1.NodeList

	if err = workloadClient.List(ctx, &nodes, client.MatchingLabels{"metal.sidero.dev/uuid": metalMachine.Spec.ServerRef.Name}); err != nil {
		return err
	}

	if len(nodes.Items) == 0 {
		return fmt.Errorf("no matching nodes found")
	}

	if len(nodes.Items) > 1 {
		return fmt.Errorf("multiple nodes found with same uuid label")
	}

	providerID := fmt.Sprintf("%s://%s", constants.ProviderID, metalMachine.Spec.ServerRef.Name)

	node := nodes.Items[0]

	if node.Spec.ProviderID == providerID {
		return nil
	}

	patchHelper, err := patch.NewHelper(&node, workloadClient)
	if err != nil {
		return err
	}

	r.Log.Info("Setting provider ID", "id", providerID)

	node.Spec.ProviderID = providerID

	return patchHelper.Patch(ctx, &node)
}

// createServerBinding updates a server to mark it as "in use" via ServerBinding resource.
func (r *MetalMachineReconciler) createServerBinding(ctx context.Context, serverClassRef *corev1.ObjectReference, serverObj *metalv1.Server, metalMachine *infrav1.MetalMachine) error {
	var serverBinding infrav1.ServerBinding

	serverBinding.Namespace = serverObj.Namespace
	serverBinding.Name = serverObj.Name
	serverBinding.Labels = make(map[string]string)
	serverBinding.Spec.MetalMachineRef = corev1.ObjectReference{
		Kind:      metalMachine.Kind,
		UID:       metalMachine.UID,
		Namespace: metalMachine.Namespace,
		Name:      metalMachine.Name,
	}

	serverBinding.Spec.ServerClassRef = serverClassRef.DeepCopy()

	for label, value := range metalMachine.Labels {
		serverBinding.Labels[label] = value
	}

	return r.Create(ctx, &serverBinding)
}

func (r *MetalMachineReconciler) fetchServerClass(ctx context.Context, classRef *corev1.ObjectReference) (*metalv1.ServerClass, error) {
	serverClassResource := &metalv1.ServerClass{}

	namespacedName := types.NamespacedName{
		Namespace: classRef.Namespace,
		Name:      classRef.Name,
	}

	if err := r.Get(ctx, namespacedName, serverClassResource); err != nil {
		return nil, err
	}

	return serverClassResource, nil
}

func (r *MetalMachineReconciler) ResetServer(ctx context.Context, metalMachine *infrav1.MetalMachine, serverBinding *infrav1.ServerBinding) error {
	var (
		talosSecret corev1.Secret
		serverObj   metalv1.Server
	)

	if err := r.Get(ctx, types.NamespacedName{Namespace: "", Name: metalMachine.Spec.ServerRef.Name}, &serverObj); err != nil {
		return err
	}

	if serverObj.Spec.BMC != nil {
		// let BMC configuration to reboot and wipe te machine using pxe.
		return nil
	}

	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, metalMachine.ObjectMeta)
	if err != nil {
		return fmt.Errorf("no cluster label or cluster does not exist")
	}

	if err = r.Get(ctx, types.NamespacedName{Namespace: cluster.Namespace, Name: fmt.Sprintf("%s-talosconfig", cluster.Name)}, &talosSecret); err != nil {
		return err
	}

	config, ok := talosSecret.Data["talosconfig"]
	if !ok {
		return fmt.Errorf("failed to find talosconfig data in the talosconfig secret")
	}

	var clientConfig *clientconfig.Config
	clientConfig, err = clientconfig.FromBytes(config)

	if err != nil {
		return err
	}

	var talosClient *talosclient.Client
	talosClient, err = talosclient.New(ctx,
		talosclient.WithConfig(clientConfig),
		talosclient.WithEndpoints(serverBinding.Spec.Addresses...),
	)
	if err != nil {
		return err
	}

	// ignore error if the machine is already reset, reboot, offline, but record event
	if err = talosClient.Reset(ctx, false, true); err != nil {
		r.Recorder.Event(metalMachine.Spec.ServerRef, corev1.EventTypeWarning, "Server Allocation", fmt.Sprintf("Software reset failed on %q, %s", metalMachine.Name, err))
	} else {
		r.Recorder.Event(metalMachine.Spec.ServerRef, corev1.EventTypeNormal, "Server Allocation", fmt.Sprintf("Software reset called on %q", metalMachine.Name))
	}

	return nil
}
