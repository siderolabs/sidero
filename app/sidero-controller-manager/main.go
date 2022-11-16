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
	capi "sigs.k8s.io/cluster-api/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	infrav1alpha3 "github.com/talos-systems/sidero/app/caps-controller-manager/api/v1alpha3"
	metalv1alpha1 "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha1"
	metalv1alpha2 "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha2"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/controllers"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/ipxe"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/metadata"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/power/api"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/server"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/siderolink"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/tftp"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/pkg/constants"
	siderotypes "github.com/talos-systems/sidero/app/sidero-controller-manager/pkg/types"
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
	_ = metalv1alpha2.AddToScheme(scheme)
	_ = infrav1alpha3.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

//nolint:maintidx
func main() {
	var (
		metricsAddr          string
		healthAddr           string
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
		ipmiPXEMethod        string

		testPowerSimulatedExplicitFailureProb float64
		testPowerSimulatedSilentFailureProb   float64
	)

	flag.StringVar(&apiEndpoint, "api-endpoint", "", "The endpoint (hostname or IP address) Sidero can be reached at from the servers.")
	flag.IntVar(&apiPort, "api-port", 8081, "The TCP port Sidero components can be reached at from the servers.")
	flag.IntVar(&httpPort, "http-port", 8081, "The TCP port Sidero controller manager HTTP server is running.")
	flag.StringVar(&metricsAddr, "metrics-bind-addr", ":8081", "The address the metric endpoint binds to.")
	flag.StringVar(&healthAddr, "health-addr", ":9440", "The address the health endpoint binds to.")
	flag.StringVar(&extraAgentKernelArgs, "extra-agent-kernel-args", "", "A list of Linux kernel command line arguments to add to the agent environment kernel parameters (e.g. 'console=tty1 console=ttyS1').")
	flag.StringVar(&bootFromDiskMethod, "boot-from-disk-method", string(siderotypes.BootIPXEExit), "Default method to use to boot server from disk if it hits iPXE endpoint after install.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", true, "Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&autoAcceptServers, "auto-accept-servers", false, "Add servers as 'accepted' when they register with Sidero API.")
	flag.BoolVar(&insecureWipe, "insecure-wipe", true, "Wipe head of the disk only (if false, wipe whole disk).")
	flag.BoolVar(&autoBMCSetup, "auto-bmc-setup", true, "Attempt to setup BMC info automatically when agent boots.")
	flag.DurationVar(&serverRebootTimeout, "server-reboot-timeout", constants.DefaultServerRebootTimeout, "Timeout to wait for the server to restart and start wipe.")
	flag.StringVar(&ipmiPXEMethod, "ipmi-pxe-method", string(siderotypes.PXEModeUEFI), fmt.Sprintf("Default method to use to set server to boot from PXE via IPMI: %s.", []string{siderotypes.PXEModeUEFI, siderotypes.PXEModeBIOS}))
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

	if !siderotypes.PXEMode(ipmiPXEMethod).IsValid() {
		setupLog.Error(fmt.Errorf("ipmi-pxe-method is invalid"), "")
		os.Exit(1)
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
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "controller-leader-election-sidero-controller-manager",
		Port:                   9443,
		HealthProbeBindAddress: healthAddr,
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

	ctx := ctrl.SetupSignalHandler()

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
		PXEMode:       siderotypes.PXEMode(ipmiPXEMethod),
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

	setupWebhooks(mgr)
	setupChecks(mgr, httpPort)

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

	if err := ipxe.RegisterIPXE(httpMux, apiEndpoint, apiPort, extraAgentKernelArgs, siderotypes.BootFromDisk(bootFromDiskMethod), apiPort, mgr.GetClient()); err != nil {
		setupLog.Error(err, "unable to start iPXE server", "controller", "Environment")
		os.Exit(1)
	}

	setupLog.Info("starting metadata server")

	if err := metadata.RegisterServer(httpMux, mgr.GetClient()); err != nil {
		setupLog.Error(err, "unable to start metadata server", "controller", "Environment")
		os.Exit(1)
	}

	setupLog.Info("starting internal API server")

	apiRecorder := eventBroadcaster.NewRecorder(
		mgr.GetScheme(),
		corev1.EventSource{Component: "sidero-server"})

	grpcServer := server.CreateServer(mgr.GetClient(), apiRecorder, mgr.GetScheme(), autoAcceptServers, insecureWipe, autoBMCSetup, serverRebootTimeout)

	if err = mgr.Add(manager.RunnableFunc(func(ctx context.Context) error {
		return siderolink.Cfg.LoadOrCreate(ctx, mgr.GetClient())
	})); err != nil {
		setupLog.Error(err, `failed to add SideroLink configuration initialization`)
		os.Exit(1)
	}

	if err = mgr.Add(RunnableClientFunc(controllers.ReconcileServerClassAny)); err != nil {
		setupLog.Error(err, `failed to add initial reconcile`)
		os.Exit(1)
	}

	if err = mgr.Add(RunnableClientFunc(func(ctx context.Context, k8sClient client.Client) error {
		return controllers.ReconcileEnvironmentDefault(ctx, k8sClient, TalosRelease, apiEndpoint, uint16(apiPort))
	})); err != nil {
		setupLog.Error(err, `failed to add initial reconcile`)
		os.Exit(1)
	}

	setupLog.Info("starting manager and HTTP server")

	go func() {
		err := mgr.Start(ctx)
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

func setupWebhooks(mgr ctrl.Manager) {
	if err := (&metalv1alpha1.ServerClass{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "ServerClass")
		os.Exit(1)
	}

	if err := (&metalv1alpha1.Environment{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "Environment")
		os.Exit(1)
	}

	if err := (&metalv1alpha1.Server{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "Server")
		os.Exit(1)
	}

	if err := (&metalv1alpha2.ServerClass{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "ServerClass")
		os.Exit(1)
	}

	if err := (&metalv1alpha2.Environment{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "Environment")
		os.Exit(1)
	}

	if err := (&metalv1alpha2.Server{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "Server")
		os.Exit(1)
	}
}

func setupChecks(mgr ctrl.Manager, httpPort int) {
	addr := fmt.Sprintf("127.0.0.1:%d", httpPort)

	if err := mgr.AddReadyzCheck("ipxe", ipxe.Check(addr)); err != nil {
		setupLog.Error(err, "unable to create ready check")
		os.Exit(1)
	}

	if err := mgr.AddHealthzCheck("webhook", ipxe.Check(addr)); err != nil {
		setupLog.Error(err, "unable to create health check")
		os.Exit(1)
	}
}

// RunnableClientFunc implements Runnable and inject.Client using a function.
func RunnableClientFunc(f func(context.Context, client.Client) error) *runnableClientFunc {
	return &runnableClientFunc{
		Func: f,
	}
}

type runnableClientFunc struct {
	Client client.Client
	Func   func(context.Context, client.Client) error
}

// InjectClient implements inject.Client.
//
//nolint:unparam
func (r *runnableClientFunc) InjectClient(c client.Client) error {
	r.Client = c

	return nil
}

// Start implements Runnable.
func (r *runnableClientFunc) Start(ctx context.Context) error {
	return r.Func(ctx, r.Client)
}
