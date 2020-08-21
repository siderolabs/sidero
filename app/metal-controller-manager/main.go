// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	metalv1alpha1 "github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/app/metal-controller-manager/controllers"
	"github.com/talos-systems/sidero/app/metal-controller-manager/internal/ipxe"
	"github.com/talos-systems/sidero/app/metal-controller-manager/internal/server"
	"github.com/talos-systems/sidero/app/metal-controller-manager/internal/tftp"
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
	// +kubebuilder:scaffold:scheme
}

func main() {
	var (
		metricsAddr          string
		apiEndpoint          string
		enableLeaderElection bool
	)

	flag.StringVar(&apiEndpoint, "api-endpoint", "", "The endpoint used by the discovery environment.")
	flag.StringVar(&metricsAddr, "metrics-addr", ":8081", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false, "Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")

	flag.Parse()

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

	if err = (&controllers.EnvironmentReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Environment"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr, controller.Options{MaxConcurrentReconciles: defaultMaxConcurrentReconciles}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Environment")
		os.Exit(1)
	}

	if err = (&controllers.ServerReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Server"),
		Scheme: mgr.GetScheme(),
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
		}
	}()

	setupLog.Info("starting iPXE server")

	go func() {
		if apiEndpoint == "" {
			if endpoint, ok := os.LookupEnv("API_ENDPOINT"); ok {
				apiEndpoint = endpoint
			} else {
				setupLog.Error(fmt.Errorf("no api endpoint found"), "unable to start iPXE server", "controller", "Environment")
			}
		}

		if err := ipxe.ServeIPXE(apiEndpoint); err != nil {
			setupLog.Error(err, "unable to start iPXE server", "controller", "Environment")
		}
	}()

	setupLog.Info("starting internal API server")

	go func() {
		if err := server.Serve(); err != nil {
			setupLog.Error(err, "unable to start API server", "controller", "Environment")
		}
	}()

	setupLog.Info("starting manager")

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
