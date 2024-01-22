// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"flag"
	"os"

	debug "github.com/siderolabs/go-debug"
	"github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	cgrecord "k8s.io/client-go/tools/record"
	"k8s.io/component-base/logs"
	logsv1 "k8s.io/component-base/logs/api/v1"
	"k8s.io/klog/v2"
	capiv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/remote"
	"sigs.k8s.io/cluster-api/util/flags"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	infrav1alpha2 "github.com/siderolabs/sidero/app/caps-controller-manager/api/v1alpha2"
	infrav1alpha3 "github.com/siderolabs/sidero/app/caps-controller-manager/api/v1alpha3"
	"github.com/siderolabs/sidero/app/caps-controller-manager/controllers"
	metalv1 "github.com/siderolabs/sidero/app/sidero-controller-manager/api/v1alpha2"
	// +kubebuilder:scaffold:imports
)

const (
	defaultMaxConcurrentReconciles = 10
	debugAddr                      = ":9994"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

//nolint:wsl
func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = capiv1.AddToScheme(scheme)
	_ = infrav1alpha2.AddToScheme(scheme)
	_ = infrav1alpha3.AddToScheme(scheme)
	_ = metalv1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

var (
	healthAddr           string
	enableLeaderElection bool
	webhookPort          int
	webhookCertDir       string
	tlsOptions           = flags.TLSOptions{}
	diagnosticsOptions   = flags.DiagnosticsOptions{}
	logOptions           = logs.NewOptions()
)

// InitFlags initializes the flags.
func InitFlags(fs *pflag.FlagSet) {
	logsv1.AddFlags(logOptions, fs)

	fs.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	fs.IntVar(&webhookPort, "webhook-port", 9443, "Webhook Server port, disabled by default. When enabled, the manager will only work as webhook server, no reconcilers are installed.")
	fs.StringVar(&webhookCertDir, "webhook-cert-dir", "/tmp/k8s-webhook-server/serving-certs/",
		"Webhook cert dir, only used when webhook-port is specified.")
	fs.StringVar(&healthAddr, "health-addr", ":9440",
		"The address the health endpoint binds to.")

	flags.AddDiagnosticsOptions(fs, &diagnosticsOptions)
	flags.AddTLSOptions(fs, &tlsOptions)
}

func main() {
	InitFlags(pflag.CommandLine)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	if err := logsv1.ValidateAndApply(logOptions, nil); err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// klog.Background will automatically use the right logger.
	ctrl.SetLogger(klog.Background())

	diagnosticsOpts := flags.GetDiagnosticsOptions(diagnosticsOptions)

	tlsOptionOverrides, err := flags.GetTLSOptionOverrideFuncs(tlsOptions)
	if err != nil {
		setupLog.Error(err, "unable to add TLS settings to the webhook server")
		os.Exit(1)
	}

	go func() {
		debugLogFunc := func(msg string) {
			setupLog.Info(msg)
		}
		if err := debug.ListenAndServe(context.TODO(), debugAddr, debugLogFunc); err != nil {
			setupLog.Error(err, "failed to start debug server")
			os.Exit(1)
		}
	}()

	// Machine and cluster operations can create enough events to trigger the event recorder spam filter
	// Setting the burst size higher ensures all events will be recorded and submitted to the API
	broadcaster := cgrecord.NewBroadcasterWithCorrelatorOptions(cgrecord.CorrelatorOptions{
		BurstSize: 100,
	})

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		Metrics:                diagnosticsOpts,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "controller-leader-election-capm",
		EventBroadcaster:       broadcaster,
		HealthProbeBindAddress: healthAddr,
		Client: client.Options{
			Cache: &client.CacheOptions{
				DisableFor: []client.Object{
					&corev1.ConfigMap{},
					&corev1.Secret{},
				},
			},
		},
		WebhookServer: webhook.NewServer(
			webhook.Options{
				Port:    webhookPort,
				CertDir: webhookCertDir,
				TLSOpts: tlsOptionOverrides,
			},
		),
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	ctx := ctrl.SetupSignalHandler()

	// Set up a ClusterCacheTracker to provide to controllers
	// requiring a connection to a remote cluster
	tracker, err := remote.NewClusterCacheTracker(mgr, remote.ClusterCacheTrackerOptions{
		Indexes:               []remote.Index{remote.NodeProviderIDIndex},
		ClientUncachedObjects: []client.Object{},
	})
	if err != nil {
		setupLog.Error(err, "unable to create cluster cache tracker")
		os.Exit(1)
	}

	if err := (&remote.ClusterCacheReconciler{
		Client:  mgr.GetClient(),
		Tracker: tracker,
	}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: defaultMaxConcurrentReconciles}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ClusterCacheReconciler")
		os.Exit(1)
	}

	recorder := mgr.GetEventRecorderFor("caps-controller-manager")

	if err = (&controllers.MetalClusterReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("MetalCluster"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: defaultMaxConcurrentReconciles}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "MetalCluster")
		os.Exit(1)
	}

	if err = (&controllers.MetalMachineReconciler{
		Client:   mgr.GetClient(),
		Log:      ctrl.Log.WithName("controllers").WithName("MetalMachine"),
		Scheme:   mgr.GetScheme(),
		Recorder: recorder,
		Tracker:  tracker,
	}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: defaultMaxConcurrentReconciles}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "MetalMachine")
		os.Exit(1)
	}

	if err = (&controllers.ServerBindingReconciler{
		Client:   mgr.GetClient(),
		Log:      ctrl.Log.WithName("controllers").WithName("ServerBinding"),
		Scheme:   mgr.GetScheme(),
		Recorder: recorder,
	}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: defaultMaxConcurrentReconciles}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ServerBinding")
		os.Exit(1)
	}

	// +kubebuilder:scaffold:builder

	setupWebhooks(mgr)
	setupChecks(mgr)

	setupLog.Info("starting manager")

	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func setupWebhooks(mgr ctrl.Manager) {
	var err error

	if err = (&infrav1alpha3.MetalCluster{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "MetalCluster")
		os.Exit(1)
	}

	if err = (&infrav1alpha3.MetalMachine{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "MetalMachine")
		os.Exit(1)
	}

	if err = (&infrav1alpha3.MetalMachineTemplate{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "MetalMachineTemplate")
		os.Exit(1)
	}

	if err = (&infrav1alpha3.ServerBinding{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "ServerBinding")
		os.Exit(1)
	}
}

func setupChecks(mgr ctrl.Manager) {
	if err := mgr.AddReadyzCheck("webhook", mgr.GetWebhookServer().StartedChecker()); err != nil {
		setupLog.Error(err, "unable to create ready check")
		os.Exit(1)
	}

	if err := mgr.AddHealthzCheck("webhook", mgr.GetWebhookServer().StartedChecker()); err != nil {
		setupLog.Error(err, "unable to create health check")
		os.Exit(1)
	}
}
