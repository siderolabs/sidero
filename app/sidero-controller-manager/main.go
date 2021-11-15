// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	debug "github.com/talos-systems/go-debug"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/record"
	capi "sigs.k8s.io/cluster-api/api/v1alpha4"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	infrav1 "github.com/talos-systems/sidero/app/caps-controller-manager/api/v1alpha3"
	metalv1alpha1 "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/controllers"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/healthz"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/ipxe"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/metadata"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/power/api"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/server"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/tftp"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/pkg/constants"
	"github.com/talos-systems/sidero/internal/client"
	// +kubebuilder:scaffold:imports
)

const (
	defaultMaxConcurrentReconciles = 10
	debugAddr                      = ":9992"
)

var (
	// TalosRelease is set as a build argument.
	TalosRelease string

	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

//nolint:wsl
func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = capi.AddToScheme(scheme)

	_ = metalv1alpha1.AddToScheme(scheme)
	_ = infrav1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var (
		metricsAddr          string
		apiEndpoint          string
		apiPort              int
		httpPort             int
		extraAgentKernelArgs string
		bootFromDiskMethod   string
		enableLeaderElection bool
		autoAcceptServers    bool
		insecureWipe         bool
		autoBMCSetup         bool
		serverRebootTimeout  time.Duration

		testPowerSimulatedExplicitFailureProb float64
		testPowerSimulatedSilentFailureProb   float64
	)

	flag.StringVar(&apiEndpoint, "api-endpoint", "", "The endpoint (hostname or IP address) Sidero can be reached at from the servers.")
	flag.IntVar(&apiPort, "api-port", 8081, "The TCP port Sidero components can be reached at from the servers.")
	flag.IntVar(&httpPort, "http-port", 8081, "The TCP port Sidero controller manager HTTP server is running.")
	flag.StringVar(&metricsAddr, "metrics-bind-addr", ":8081", "The address the metric endpoint binds to.")
	flag.StringVar(&extraAgentKernelArgs, "extra-agent-kernel-args", "", "A list of Linux kernel command line arguments to add to the agent environment kernel parameters (e.g. 'console=tty1 console=ttyS1').")
	flag.StringVar(&bootFromDiskMethod, "boot-from-disk-method", string(ipxe.BootIPXEExit), "Default method to use to boot server from disk if it hits iPXE endpoint after install.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", true, "Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&autoAcceptServers, "auto-accept-servers", false, "Add servers as 'accepted' when they register with Sidero API.")
	flag.BoolVar(&insecureWipe, "insecure-wipe", true, "Wipe head of the disk only (if false, wipe whole disk).")
	flag.BoolVar(&autoBMCSetup, "auto-bmc-setup", true, "Attempt to setup BMC info automatically when agent boots.")
	flag.DurationVar(&serverRebootTimeout, "server-reboot-timeout", constants.DefaultServerRebootTimeout, "Timeout to wait for the server to restart and start wipe.")
	flag.Float64Var(&testPowerSimulatedExplicitFailureProb, "test-power-simulated-explicit-failure-prob", 0, "Test failure simulation setting.")
	flag.Float64Var(&testPowerSimulatedSilentFailureProb, "test-power-simulated-silent-failure-prob", 0, "Test failure simulation setting.")

	flag.Parse()

	// we can't continue without it
	if TalosRelease == "" {
		panic("TalosRelease is not set during the build")
	}

	// workaround for clusterctl not accepting empty value as default value
	if extraAgentKernelArgs == "-" {
		extraAgentKernelArgs = ""
	}

	if apiEndpoint == "-" {
		apiEndpoint = ""
	}

	if apiEndpoint == "" {
		if endpoint, ok := os.LookupEnv("API_ENDPOINT"); ok {
			apiEndpoint = endpoint
		} else {
			setupLog.Error(fmt.Errorf("no api endpoint found"), "")
			os.Exit(1)
		}
	}

	ctrl.SetLogger(zap.New(func(o *zap.Options) {
		o.Development = true
	}))

	go func() {
		debugLogFunc := func(msg string) {
			setupLog.Info(msg)
		}
		if err := debug.ListenAndServe(context.TODO(), debugAddr, debugLogFunc); err != nil {
			setupLog.Error(err, "failed to start debug server")
			os.Exit(1)
		}
	}()

	// only for testing, doesn't affect production, default values simulate no failures
	api.DefaultDice = api.NewFailureDice(testPowerSimulatedExplicitFailureProb, testPowerSimulatedSilentFailureProb)

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "controller-leader-election-sidero-controller-manager",
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

	ctx := context.Background()

	if err = (&controllers.EnvironmentReconciler{
		Client:       mgr.GetClient(),
		Log:          ctrl.Log.WithName("controllers").WithName("Environment"),
		Scheme:       mgr.GetScheme(),
		TalosRelease: TalosRelease,
		APIEndpoint:  apiEndpoint,
		APIPort:      uint16(apiPort),
	}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: defaultMaxConcurrentReconciles}); err != nil {
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
	}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: defaultMaxConcurrentReconciles}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Server")
		os.Exit(1)
	}

	if err = (&controllers.ServerClassReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("ServerClass"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: defaultMaxConcurrentReconciles}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ServerClass")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	errCh := make(chan error)

	setupLog.Info("starting TFTP server")

	go func() {
		if err := tftp.ServeTFTP(); err != nil {
			setupLog.Error(err, "unable to start TFTP server", "controller", "Environment")
		}

		errCh <- err
	}()

	httpMux := http.NewServeMux()

	setupLog.Info("starting iPXE server")

	if err := ipxe.RegisterIPXE(httpMux, apiEndpoint, apiPort, extraAgentKernelArgs, ipxe.BootFromDisk(bootFromDiskMethod), apiPort, mgr.GetClient()); err != nil {
		setupLog.Error(err, "unable to start iPXE server", "controller", "Environment")
		os.Exit(1)
	}

	setupLog.Info("starting metadata server")

	if err := metadata.RegisterServer(httpMux, mgr.GetClient()); err != nil {
		setupLog.Error(err, "unable to start metadata server", "controller", "Environment")
		os.Exit(1)
	}

	setupLog.Info("starting healthz server")

	if err := healthz.RegisterServer(httpMux); err != nil {
		setupLog.Error(err, "unable to start healthz server", "controller", "Environment")
		os.Exit(1)
	}

	setupLog.Info("starting internal API server")

	apiRecorder := eventBroadcaster.NewRecorder(
		mgr.GetScheme(),
		corev1.EventSource{Component: "sidero-server"})

	grpcServer := server.CreateServer(mgr.GetClient(), apiRecorder, mgr.GetScheme(), autoAcceptServers, insecureWipe, autoBMCSetup, serverRebootTimeout)

	k8sClient, err := client.NewClient(nil)
	if err != nil {
		setupLog.Error(err, `failed to create k8s client`)
		os.Exit(1)
	}

	if err = controllers.ReconcileServerClassAny(ctx, k8sClient); err != nil {
		setupLog.Error(err, `failed to reconcile ServerClass "any"`)
		os.Exit(1)
	}

	if err = controllers.ReconcileEnvironmentDefault(ctx, k8sClient, TalosRelease, apiEndpoint, uint16(apiPort)); err != nil {
		setupLog.Error(err, `failed to reconcile Environment "default"`)
		os.Exit(1)
	}

	setupLog.Info("starting manager and HTTP server")

	go func() {
		err := mgr.Start(ctrl.SetupSignalHandler())
		if err != nil {
			setupLog.Error(err, "problem running manager")
		}

		errCh <- err
	}()

	go func() {
		// Go standard library doesn't support running HTTP/2 on non-TLS HTTP connections.
		// Package h2c provides handling for HTTP/2 over plaintext connection.
		// gRPC provides its own HTTP/2 server implementation, so that's not an issue for gRPC,
		// but as we unify all endpoints under a single HTTP endpoint, we have to provide additional
		// layer of support here.
		h2s := &http2.Server{}

		grpcHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.ProtoMajor == 2 && strings.HasPrefix(
				req.Header.Get("Content-Type"), "application/grpc") {
				// grpcServer provides internal gRPC API server
				grpcServer.ServeHTTP(w, req)

				return
			}

			// httpMux contains iPXE server and metadata server handlers
			httpMux.ServeHTTP(w, req)
		})

		err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), h2c.NewHandler(grpcHandler, h2s))
		if err != nil {
			setupLog.Error(err, "problem running HTTP server")
		}

		errCh <- err
	}()

	for err = range errCh {
		if err != nil {
			os.Exit(1)
		}
	}
}
