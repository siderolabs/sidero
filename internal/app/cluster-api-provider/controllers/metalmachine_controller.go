// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	infrav1 "github.com/talos-systems/sidero/internal/app/cluster-api-provider/api/v1alpha3"
	"github.com/talos-systems/sidero/internal/app/cluster-api-provider/internal/pkg/ipmi"
	"github.com/talos-systems/sidero/internal/app/cluster-api-provider/pkg/constants"
	metalv1alpha1 "github.com/talos-systems/sidero/internal/app/metal-controller-manager/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/pointer"
	capiv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// MetalMachineReconciler reconciles a MetalMachine object
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

func (r *MetalMachineReconciler) Reconcile(req ctrl.Request) (_ ctrl.Result, rerr error) {
	ctx := context.Background()
	logger := r.Log.WithValues("metalmachine", req.NamespacedName)

	// Fetch the metalMachine instance.
	metalMachine := &infrav1.MetalMachine{}
	err := r.Get(ctx, req.NamespacedName, metalMachine)
	if apierrors.IsNotFound(err) {
		return ctrl.Result{}, nil
	}
	if err != nil {
		return ctrl.Result{}, err
	}

	machine, err := util.GetOwnerMachine(ctx, r.Client, metalMachine.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if machine == nil {
		return ctrl.Result{}, fmt.Errorf("no ownerref for machine %s", metalMachine.ObjectMeta.Name)
	}

	logger = logger.WithName(fmt.Sprintf("machine=%s", machine.Name))

	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, machine.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("no cluster label or cluster does not exist")
	}

	logger = logger.WithName(fmt.Sprintf("cluster=%s", cluster.Name))

	if !cluster.Status.InfrastructureReady {
		return ctrl.Result{}, fmt.Errorf("clusterInfra is not available yet")
	}

	if machine.Spec.Bootstrap.DataSecretName == nil {
		return ctrl.Result{}, fmt.Errorf("bootstrap secret is not available yet")
	}

	// Initialize the patch helper
	patchHelper, err := patch.NewHelper(metalMachine, r)
	if err != nil {
		return ctrl.Result{}, err
	}
	// Always attempt to Patch the argesMachine object and status after each reconciliation.
	defer func() {
		if err := patchHelper.Patch(ctx, metalMachine); err != nil {
			logger.Error(err, "failed to patch metalMachine")
			if rerr == nil {
				rerr = err
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
			Namespace: metalMachine.Spec.ServerRef.Namespace,
			Name:      metalMachine.Spec.ServerRef.Name,
		}

		r.Get(ctx, namespacedName, serverResource)
	} else {
		if metalMachine.Spec.ServerClassRef == nil {
			return ctrl.Result{}, fmt.Errorf("either a server or serverclass ref must be supplied")
		}

		serverResource, err = r.fetchServerFromClass(ctx, metalMachine.Spec.ServerClassRef)
		if err != nil {
			return ctrl.Result{}, err
		}

		metalMachine.Spec.ServerRef = &v1.ObjectReference{
			Kind: serverResource.Kind,
			Name: serverResource.Name,
		}
	}

	if serverResource.Spec.BMC != nil {
		ipmiClient, err := ipmi.NewClient(*serverResource.Spec.BMC)
		if err != nil {
			return ctrl.Result{}, err
		}

		chassisStatus, err := ipmiClient.Status()
		if err != nil {
			return ctrl.Result{}, err
		}

		// Only take action if server is turned off
		// otherwise IPMI library gets angry
		if !chassisStatus.IsSystemPowerOn() {
			err = ipmiClient.SetPXE()
			if err != nil {
				return ctrl.Result{}, err
			}

			err = ipmiClient.PowerOn()
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	// Set the providerID, as its required in upstream capi for machine lifecycle
	metalMachine.Spec.ProviderID = pointer.StringPtr(fmt.Sprintf("%s://%s", constants.ProviderID, serverResource.Name))

	err = r.patchProviderID(ctx, cluster, metalMachine)
	if err != nil {
		return ctrl.Result{}, err
	}

	metalMachine.Status.Ready = true
	return ctrl.Result{}, nil
}

func (r *MetalMachineReconciler) reconcileDelete(ctx context.Context, metalMachine *infrav1.MetalMachine) (ctrl.Result, error) {
	//TODO(rsmitty): add in call to reset node via talos api

	serverResource := &metalv1alpha1.Server{}

	// Use server ref if already provided
	if metalMachine.Spec.ServerRef != nil {
		namespacedName := types.NamespacedName{
			Namespace: metalMachine.Spec.ServerRef.Namespace,
			Name:      metalMachine.Spec.ServerRef.Name,
		}

		r.Get(ctx, namespacedName, serverResource)
	}

	if serverResource.Spec.BMC != nil {
		ipmiClient, err := ipmi.NewClient(*serverResource.Spec.BMC)
		if err != nil {
			return ctrl.Result{}, err
		}

		err = ipmiClient.SetPXE()
		if err != nil {
			return ctrl.Result{}, err
		}

		err = ipmiClient.PowerOff()
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	// Remove in-use bool from server object
	if metalMachine.Spec.ServerRef != nil {
		serverObj := &metalv1alpha1.Server{}

		namespacedName := types.NamespacedName{
			Namespace: "",
			Name:      metalMachine.Spec.ServerRef.Name,
		}

		if err := r.Get(ctx, namespacedName, serverObj); err != nil {
			return ctrl.Result{}, err
		}

		serverObj.Status.InUse = false

		if err := r.Status().Update(ctx, serverObj); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Cluster is deleted so remove the finalizer.
	controllerutil.RemoveFinalizer(metalMachine, infrav1.MachineFinalizer)

	return ctrl.Result{}, nil
}

func (r *MetalMachineReconciler) SetupWithManager(mgr ctrl.Manager, options controller.Options) error {
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1.MetalMachine{}).
		Complete(r)
}

func (r *MetalMachineReconciler) fetchServerFromClass(ctx context.Context, classRef *corev1.ObjectReference) (*metalv1alpha1.Server, error) {
	// Grab server class and validate that we have nodes available
	serverClassResource := &metalv1alpha1.ServerClass{}

	namespacedName := types.NamespacedName{
		Namespace: classRef.Namespace,
		Name:      classRef.Name,
	}

	if err := r.Get(ctx, namespacedName, serverClassResource); err != nil {
		return nil, err
	}

	if len(serverClassResource.Status.ServersAvailable) == 0 {
		return nil, fmt.Errorf("no servers available in serverclass")
	}

	// Fetch first server in avail list
	firstAvailServer := serverClassResource.Status.ServersAvailable[0]
	serverObj := &metalv1alpha1.Server{}

	namespacedName = types.NamespacedName{
		Namespace: "",
		Name:      firstAvailServer,
	}

	if err := r.Get(ctx, namespacedName, serverObj); err != nil {
		return nil, err
	}

	// Update server status to in use
	serverObj.Status.InUse = true

	if err := r.Status().Update(ctx, serverObj); err != nil {
		return nil, err
	}

	return serverObj, nil
}

func (r *MetalMachineReconciler) patchProviderID(ctx context.Context, cluster *capiv1.Cluster, metalMachine *infrav1.MetalMachine) error {
	kubeconfigSecret := &v1.Secret{}
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

	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	nodes, err = clientset.CoreV1().Nodes().List(
		ctx,
		metav1.ListOptions{
			LabelSelector: fmt.Sprintf("metal.sidero.dev/uuid=%s", metalMachine.Spec.ServerRef.Name),
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

	for _, node := range nodes.Items {
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
