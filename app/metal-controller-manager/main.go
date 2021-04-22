// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	infrav1 "github.com/talos-systems/sidero/app/cluster-api-provider-sidero/api/v1alpha3"
	metalv1alpha1 "github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/app/metal-controller-manager/controllers"
	"github.com/talos-systems/sidero/app/metal-controller-manager/internal/ipxe"
	"github.com/talos-systems/sidero/app/metal-controller-manager/internal/power/api"
	"github.com/talos-systems/sidero/app/metal-controller-manager/internal/server"
	"github.com/talos-systems/sidero/app/metal-controller-manager/internal/tftp"
	"github.com/talos-systems/sidero/app/metal-controller-manager/pkg/constants"
	// +kubebuilder:scaffold:imports
)

const (
	defaultMaxConcurrentReconciles = 10
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

// nolint: wsl
func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = metalv1alpha1.AddToScheme(scheme)
	_ = infrav1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var (
		metricsAddr          string
		apiEndpoint          string
		extraAgentKernelArgs string
		enableLeaderElection bool
		autoAcceptServers    bool
		insecureWipe         bool
		serverRebootTimeout  time.Duration

		testPowerSimulatedExplicitFailureProb float64
		testPowerSimulatedSilentFailureProb   float64
	)

	flag.StringVar(&apiEndpoint, "api-endpoint", "", "The endpoint used by the discovery environment.")
	flag.StringVar(&metricsAddr, "metrics-addr", ":8081", "The address the metric endpoint binds to.")
	flag.StringVar(&extraAgentKernelArgs, "extra-agent-kernel-args", "", "A comma delimited list of key-value pairs to be added to the agent environment kernel parameters.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", true, "Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&autoAcceptServers, "auto-accept-servers", false, "Add servers as 'accepted' when they register with Sidero API.")
	flag.BoolVar(&insecureWipe, "insecure-wipe", true, "Wipe head of the disk only (if false, wipe whole disk).")
	flag.DurationVar(&serverRebootTimeout, "server-reboot-timeout", constants.DefaultServerRebootTimeout, "Timeout to wait for the server to restart and start wipe.")
	flag.Float64Var(&testPowerSimulatedExplicitFailureProb, "test-power-simulated-explicit-failure-prob", 0, "Test failure simulation setting.")
	flag.Float64Var(&testPowerSimulatedSilentFailureProb, "test-power-simulated-silent-failure-prob", 0, "Test failure simulation setting.")

	flag.Parse()

	// workaround for clusterctl not accepting empty value as default value
	if extraAgentKernelArgs == "-" {
		extraAgentKernelArgs = ""
	}

	if apiEndpoint == "-" {
		apiEndpoint = ""
	}

	// only for testing, doesn't affect production, default values simulate no failures
	api.DefaultDice = api.NewFailureDice(testPowerSimulatedExplicitFailureProb, testPowerSimulatedSilentFailureProb)

	ctrl.SetLogger(zap.New(func(o *zap.Options) {
		o.Development = true
	}))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "controller-leader-election-metal-controller-manager",
		Port:               9443,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(mgr.GetConfig())
	if err != nil {
		setupLog.Error(err, "unable to create k8s client")
		os.Exit(1)
	}

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartRecordingToSink(
		&typedcorev1.EventSinkImpl{
			Interface: clientset.CoreV1().Events(""),
		})

	recorder := eventBroadcaster.NewRecorder(
		mgr.GetScheme(),
		corev1.EventSource{Component: "sidero-controller-manager"})

	if err = (&controllers.EnvironmentReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Environment"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr, controller.Options{MaxConcurrentReconciles: defaultMaxConcurrentReconciles}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Environment")
		os.Exit(1)
	}

	if err = (&controllers.ServerReconciler{
		Client:        mgr.GetClient(),
		Log:           ctrl.Log.WithName("controllers").WithName("Server"),
		Scheme:        mgr.GetScheme(),
		APIReader:     mgr.GetAPIReader(),
		Recorder:      recorder,
		RebootTimeout: serverRebootTimeout,
	}).SetupWithManager(mgr, controller.Options{MaxConcurrentReconciles: defaultMaxConcurrentReconciles}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Server")
		os.Exit(1)
	}

	if err = (&controllers.ServerClassReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("ServerClass"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr, controller.Options{MaxConcurrentReconciles: defaultMaxConcurrentReconciles}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ServerClass")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting TFTP server")

	go func() {
		if err := tftp.ServeTFTP(); err != nil {
			setupLog.Error(err, "unable to start TFTP server", "controller", "Environment")
			os.Exit(1)
		}
	}()

	setupLog.Info("starting iPXE server")

	go func() {
		if apiEndpoint == "" {
			if endpoint, ok := os.LookupEnv("API_ENDPOINT"); ok {
				apiEndpoint = endpoint
			} else {
				setupLog.Error(fmt.Errorf("no api endpoint found"), "unable to start iPXE server", "controller", "Environment")
				os.Exit(1)
			}
		}

		if err := ipxe.ServeIPXE(apiEndpoint, extraAgentKernelArgs, mgr.GetClient()); err != nil {
			setupLog.Error(err, "unable to start iPXE server", "controller", "Environment")
			os.Exit(1)
		}
	}()

	setupLog.Info("starting internal API server")

	go func() {
		recorder := eventBroadcaster.NewRecorder(
			mgr.GetScheme(),
			corev1.EventSource{Component: "sidero-server"})

		if err := server.Serve(mgr.GetClient(), recorder, mgr.GetScheme(), autoAcceptServers, insecureWipe, serverRebootTimeout); err != nil {
			setupLog.Error(err, "unable to start API server", "controller", "Environment")
			os.Exit(1)
		}
	}()

	setupLog.Info("starting manager")

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
