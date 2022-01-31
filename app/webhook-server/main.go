// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	infrav1alpha2 "github.com/talos-systems/sidero/app/caps-controller-manager/api/v1alpha2"
	infrav1alpha3 "github.com/talos-systems/sidero/app/caps-controller-manager/api/v1alpha3"
	metalv1alpha1 "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha1"
	metalv1alpha2 "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha2"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"log"
	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	scheme = runtime.NewScheme()
)

//nolint:wsl
func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = capiv1beta1.AddToScheme(scheme)
	_ = infrav1alpha2.AddToScheme(scheme)
	_ = infrav1alpha3.AddToScheme(scheme)
	_ = metalv1alpha1.AddToScheme(scheme)
	_ = metalv1alpha2.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var (
		metricsAddr string
		healthAddr  string
		webhookPort int
	)
	flag.StringVar(&metricsAddr, "metrics-bind-addr", ":8081", "The address the metric endpoint binds to.")
	flag.StringVar(&healthAddr, "health-addr", ":9440", "The address the health endpoint binds to.")
	flag.IntVar(&webhookPort, "webhook-port", 9443, "The TCP port the Webhook Server can be reached at.")

	flag.Parse()

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		LeaderElection:         false,
		Port:                   webhookPort,
		HealthProbeBindAddress: healthAddr,
	})
	if err != nil {
		log.Fatal(err, "unable to start manager")
	}

	ctx := ctrl.SetupSignalHandler()

	setupWebhooks(mgr)
	setupChecks(mgr)

	log.Print("starting manager")

	if err := mgr.Start(ctx); err != nil {
		log.Fatal(err, "problem running manager")
	}
}

func setupWebhooks(mgr ctrl.Manager) {
	if err := (&infrav1alpha3.MetalCluster{}).SetupWebhookWithManager(mgr); err != nil {
		log.Fatal(err, "unable to create webhook", "webhook", "MetalCluster")
	}

	if err := (&infrav1alpha3.MetalMachine{}).SetupWebhookWithManager(mgr); err != nil {
		log.Fatal(err, "unable to create webhook", "webhook", "MetalMachine")
	}

	if err := (&infrav1alpha3.MetalMachineTemplate{}).SetupWebhookWithManager(mgr); err != nil {
		log.Fatal(err, "unable to create webhook", "webhook", "MetalMachineTemplate")
	}

	if err := (&infrav1alpha3.ServerBinding{}).SetupWebhookWithManager(mgr); err != nil {
		log.Fatal(err, "unable to create webhook", "webhook", "ServerBinding")
	}

	if err := (&metalv1alpha2.Environment{}).SetupWebhookWithManager(mgr); err != nil {
		log.Fatal(err, "unable to create webhook", "webhook", "Environment")
	}

	if err := (&metalv1alpha2.Server{}).SetupWebhookWithManager(mgr); err != nil {
		log.Fatal(err, "unable to create webhook", "webhook", "Server")
	}

	if err := (&metalv1alpha2.ServerClass{}).SetupWebhookWithManager(mgr); err != nil {
		log.Fatal(err, "unable to create webhook", "webhook", "ServerClass")
	}
}

func setupChecks(mgr ctrl.Manager) {
	if err := mgr.AddReadyzCheck("webhook", mgr.GetWebhookServer().StartedChecker()); err != nil {
		log.Fatal(err, "unable to create ready check")
	}

	if err := mgr.AddHealthzCheck("webhook", mgr.GetWebhookServer().StartedChecker()); err != nil {
		log.Fatal(err, "unable to create health check")
	}
}
