// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/reference"
	"k8s.io/utils/pointer"
	capiv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	infrav1 "github.com/talos-systems/sidero/app/cluster-api-provider-sidero/api/v1alpha3"
	"github.com/talos-systems/sidero/app/cluster-api-provider-sidero/pkg/constants"
	metalv1alpha1 "github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/internal/pkg/metal"
)

var ErrNoServersInServerClass = errors.New("no servers available in serverclass")

// MetalMachineReconciler reconciles a MetalMachine object.
type MetalMachineReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=metalmachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=metalmachines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machines/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=metal.sidero.dev,resources=serverclasses,verbs=get;list;watch;
// +kubebuilder:rbac:groups=metal.sidero.dev,resources=serverclasses/status,verbs=get;list;watch;
// +kubebuilder:rbac:groups=metal.sidero.dev,resources=servers,verbs=get;list;watch;
// +kubebuilder:rbac:groups=metal.sidero.dev,resources=servers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *MetalMachineReconciler) Reconcile(req ctrl.Request) (_ ctrl.Result, err error) {
	ctx := context.Background()
	logger := r.Log.WithValues("metalmachine", req.NamespacedName)

	// Fetch the metalMachine instance.
	metalMachine := &infrav1.MetalMachine{}

	err = r.Get(ctx, req.NamespacedName, metalMachine)
	if apierrors.IsNotFound(err) {
		return ctrl.Result{RequeueAfter: constants.DefaultRequeueAfter}, nil
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
		logger.Info(" Bootstrap secret is not available yet")

		return ctrl.Result{RequeueAfter: constants.DefaultRequeueAfter}, nil
	}

	// Initialize the patch helper
	patchHelper, err := patch.NewHelper(metalMachine, r)
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

	controllerutil.AddFinalizer(metalMachine, infrav1.MachineFinalizer)

	// Handle deleted machines
	if !metalMachine.ObjectMeta.DeletionTimestamp.IsZero() {
		logger.Info("deleting machine")
		return r.reconcileDelete(ctx, metalMachine)
	}

	serverResource := &metalv1alpha1.Server{}

	// Use server ref if already provided
	if metalMachine.Spec.ServerRef != nil {
		namespacedName := types.NamespacedName{
			Namespace: "",
			Name:      metalMachine.Spec.ServerRef.Name,
		}

		if err = r.Get(ctx, namespacedName, serverResource); err != nil {
			return ctrl.Result{}, err
		}

		// Handles the case of users specifying a server ref directly and pointing to a non-accepted server
		if !serverResource.Spec.Accepted {
			return ctrl.Result{}, fmt.Errorf("specified serverref is a non-accepted server")
		}

		// Fetch the serverclass if it exists so we can ensure the server has it as an owner ref.
		var serverClassResource *metalv1alpha1.ServerClass
		if metalMachine.Spec.ServerClassRef != nil {
			serverClassResource, err = r.fetchServerClass(ctx, metalMachine.Spec.ServerClassRef)
			if err != nil {
				if errors.Is(err, ErrNoServersInServerClass) {
					return ctrl.Result{RequeueAfter: constants.DefaultRequeueAfter}, nil
				}

				return ctrl.Result{}, err
			}
		}

		// double check server is already marked in use
		// this is especially needed after pivoting the cluster from bootstrap -> mgmt plane
		if err = r.patchServerInUse(ctx, serverClassResource, serverResource, metalMachine); err != nil {
			return ctrl.Result{}, err
		}
	} else {
		if metalMachine.Spec.ServerClassRef == nil {
			return ctrl.Result{}, fmt.Errorf("either a server or serverclass ref must be supplied")
		}

		serverResource, err = r.fetchServerFromClass(ctx, metalMachine.Spec.ServerClassRef, metalMachine)
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
	}

	serverRef, err := reference.GetReference(r.Scheme, serverResource)
	if err != nil {
		return ctrl.Result{}, err
	}

	mgmtClient, err := metal.NewManagementClient(&serverResource.Spec)
	if err != nil {
		r.Recorder.Event(serverRef, corev1.EventTypeWarning, "Server Management", fmt.Sprintf("Failed to initialize management client: %s.", err))

		return ctrl.Result{}, err
	}

	poweredOn, err := mgmtClient.IsPoweredOn()
	if err != nil {
		r.Recorder.Event(serverRef, corev1.EventTypeWarning, "Server Management", fmt.Sprintf("Failed to determine power status: %s.", err))

		return ctrl.Result{}, err
	}

	// Only take action if server is turned off
	// otherwise IPMI library gets angry
	if !poweredOn {
		err = mgmtClient.SetPXE()
		if err != nil {
			r.Recorder.Event(serverRef, corev1.EventTypeWarning, "Server Management", fmt.Sprintf("Failed to set to PXE boot once: %s.", err))

			return ctrl.Result{}, err
		}

		err = mgmtClient.PowerOn()
		if err != nil {
			r.Recorder.Event(serverRef, corev1.EventTypeWarning, "Server Management", fmt.Sprintf("Failed to power on: %s.", err))

			return ctrl.Result{}, err
		}

		if !mgmtClient.IsFake() {
			r.Recorder.Event(serverRef, corev1.EventTypeNormal, "Server Management", "Server powered on.")
		}
	}

	// Set the providerID, as its required in upstream capi for machine lifecycle
	metalMachine.Spec.ProviderID = pointer.StringPtr(fmt.Sprintf("%s://%s", constants.ProviderID, serverResource.Name))

	err = r.patchProviderID(ctx, cluster, metalMachine)
	if err != nil {
		logger.Info("Failed to set provider ID", "error", err)

		return ctrl.Result{RequeueAfter: constants.DefaultRequeueAfter}, nil
	}

	metalMachine.Status.Ready = true

	return ctrl.Result{}, nil
}

func (r *MetalMachineReconciler) reconcileDelete(ctx context.Context, metalMachine *infrav1.MetalMachine) (ctrl.Result, error) {
	serverResource := &metalv1alpha1.Server{}

	if metalMachine.Spec.ServerRef != nil {
		namespacedName := types.NamespacedName{
			Namespace: "",
			Name:      metalMachine.Spec.ServerRef.Name,
		}

		if err := r.Get(ctx, namespacedName, serverResource); err != nil {
			// bail early if the server can't be fetch. this likely means we've deleted the server from underneath already.
			if apierrors.IsNotFound(err) {
				r.Log.Info("Matching server not found for metalmachine's serverref. Assuming we're orphaned.")
				controllerutil.RemoveFinalizer(metalMachine, infrav1.MachineFinalizer)

				return ctrl.Result{}, nil
			}

			return ctrl.Result{}, err
		}

		ref, err := reference.GetReference(r.Scheme, serverResource)
		if err != nil {
			return ctrl.Result{}, err
		}

		patchHelper, err := patch.NewHelper(serverResource, r)
		if err != nil {
			return ctrl.Result{}, err
		}

		serverResource.Status.InUse = false
		serverResource.OwnerReferences = []metav1.OwnerReference{}

		if err := patchHelper.Patch(ctx, serverResource); err != nil {
			return ctrl.Result{}, err
		}

		r.Recorder.Event(ref, corev1.EventTypeNormal, "Server Allocation", "Server marked as unallocated.")
	}

	// Machine is deleted so remove the finalizer.
	controllerutil.RemoveFinalizer(metalMachine, infrav1.MachineFinalizer)

	return ctrl.Result{}, nil
}

func (r *MetalMachineReconciler) SetupWithManager(mgr ctrl.Manager, options controller.Options) error {
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1.MetalMachine{}).
		Complete(r)
}

func (r *MetalMachineReconciler) fetchServerFromClass(ctx context.Context, classRef *corev1.ObjectReference, metalMachine *infrav1.MetalMachine) (*metalv1alpha1.Server, error) {
	// Grab server class and validate that we have nodes available
	serverClassResource, err := r.fetchServerClass(ctx, classRef)
	if err != nil {
		return nil, err
	}

	if len(serverClassResource.Status.ServersAvailable) == 0 {
		return nil, ErrNoServersInServerClass
	}

	// Fetch server from available list
	// NB: we added this loop to double check that an available server isn't "in use" because
	//     we saw raciness between server selection and it being removed from the ServersAvailable list.
	for _, availServer := range serverClassResource.Status.ServersAvailable {
		serverObj := &metalv1alpha1.Server{}

		namespacedName := types.NamespacedName{
			Namespace: "",
			Name:      availServer,
		}

		if err := r.Get(ctx, namespacedName, serverObj); err != nil {
			return nil, err
		}

		if serverObj.Status.InUse {
			continue
		}

		if !serverObj.Status.IsClean {
			continue
		}

		// patch server with in use bool
		if err := r.patchServerInUse(ctx, serverClassResource, serverObj, metalMachine); err != nil {
			// the server we picked was updated by another metalmachine before we finished.
			// move on to the next one.
			if apierrors.IsConflict(err) {
				continue
			}

			return nil, err
		}

		return serverObj, nil
	}

	return nil, ErrNoServersInServerClass
}

func (r *MetalMachineReconciler) patchProviderID(ctx context.Context, cluster *capiv1.Cluster, metalMachine *infrav1.MetalMachine) error {
	kubeconfigSecret := &corev1.Secret{}

	err := r.Client.Get(ctx,
		types.NamespacedName{
			Namespace: cluster.Namespace,
			Name:      cluster.Name + "-kubeconfig",
		},
		kubeconfigSecret,
	)
	if err != nil {
		return err
	}

	config, err := clientcmd.RESTConfigFromKubeConfig(kubeconfigSecret.Data["value"])
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	label := fmt.Sprintf("metal.sidero.dev/uuid=%s", metalMachine.Spec.ServerRef.Name)

	r.Log.Info("Searching for node", "label", label)

	nodes, err := clientset.CoreV1().Nodes().List(
		ctx,
		metav1.ListOptions{
			LabelSelector: label,
		},
	)
	if err != nil {
		return err
	}

	if len(nodes.Items) == 0 {
		return fmt.Errorf("no matching nodes found")
	}

	if len(nodes.Items) > 1 {
		return fmt.Errorf("multiple nodes found with same uuid label")
	}

	providerID := fmt.Sprintf("%s://%s", constants.ProviderID, metalMachine.Spec.ServerRef.Name)

	r.Log.Info("Setting provider ID", "id", providerID)

	for _, node := range nodes.Items {
		node := node

		if node.Spec.ProviderID == providerID {
			continue
		}

		node.Spec.ProviderID = providerID

		_, err = clientset.CoreV1().Nodes().Update(ctx, &node, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

// patchServerInUse updates a server to mark it as "in use".
func (r *MetalMachineReconciler) patchServerInUse(ctx context.Context, serverClass *metalv1alpha1.ServerClass, serverObj *metalv1alpha1.Server, metalMachine *infrav1.MetalMachine) error {
	if !serverObj.Status.InUse || serverObj.Status.IsClean {
		ref, err := reference.GetReference(r.Scheme, serverObj)
		if err != nil {
			return err
		}

		serverObj.Status.InUse = true
		serverObj.Status.IsClean = false

		// nb: we update status and then update the object separately b/c statuses don't seem to get
		// updated when doing the whole object below.
		if err := r.Status().Update(ctx, serverObj); err != nil {
			return err
		}

		r.Recorder.Event(ref, corev1.EventTypeNormal, "Server Allocation", fmt.Sprintf("Server marked as allocated for metalMachine %q", metalMachine.Name))
	}

	if serverClass != nil {
		patchHelper, err := patch.NewHelper(serverObj, r)
		if err != nil {
			return err
		}

		serverObj.OwnerReferences = []metav1.OwnerReference{
			*metav1.NewControllerRef(serverClass, metalv1alpha1.GroupVersion.WithKind("ServerClass")),
		}

		if err := patchHelper.Patch(ctx, serverObj); err != nil {
			return err
		}
	}

	return nil
}

func (r *MetalMachineReconciler) fetchServerClass(ctx context.Context, classRef *corev1.ObjectReference) (*metalv1alpha1.ServerClass, error) {
	serverClassResource := &metalv1alpha1.ServerClass{}

	namespacedName := types.NamespacedName{
		Namespace: classRef.Namespace,
		Name:      classRef.Name,
	}

	if err := r.Get(ctx, namespacedName, serverClassResource); err != nil {
		return nil, err
	}

	return serverClassResource, nil
}
