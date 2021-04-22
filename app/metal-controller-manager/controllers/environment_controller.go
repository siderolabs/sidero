// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-logr/logr"
	multierror "github.com/hashicorp/go-multierror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	metalv1alpha1 "github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/app/metal-controller-manager/pkg/constants"
)

// EnvironmentReconciler reconciles a Environment object.
type EnvironmentReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=metal.sidero.dev,resources=environments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal.sidero.dev,resources=environments/status,verbs=get;update;patch

func (r *EnvironmentReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	return r.reconcile(req)
}

func (r *EnvironmentReconciler) SetupWithManager(mgr ctrl.Manager, options controller.Options) error {
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&metalv1alpha1.Environment{}).
		Complete(r)
}

func (r *EnvironmentReconciler) reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()

	l := r.Log.WithValues("environment", req.Name)

	var env metalv1alpha1.Environment

	if err := r.Get(ctx, req.NamespacedName, &env); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, fmt.Errorf("unable to get environment: %w", err)
	}

	envs := filepath.Join("/var/lib/sidero/env", env.GetName())

	if _, err := os.Stat(envs); os.IsNotExist(err) {
		if err = os.MkdirAll(envs, 0o777); err != nil {
			return ctrl.Result{}, fmt.Errorf("error creating environment directory: %w", err)
		}
	}

	var (
		conditions = []metalv1alpha1.AssetCondition{}
		wg         sync.WaitGroup
		mu         sync.Mutex
		result     *multierror.Error
	)

	for _, assetTask := range []struct {
		BaseName string
		Asset    metalv1alpha1.Asset
	}{
		{
			BaseName: constants.KernelAsset,
			Asset:    env.Spec.Kernel.Asset,
		},
		{
			BaseName: constants.InitrdAsset,
			Asset:    env.Spec.Initrd.Asset,
		},
	} {
		assetTask := assetTask

		file := filepath.Join(envs, assetTask.BaseName)

		setReady := func(ready bool) {
			status := "False"
			if ready {
				status = "True"
			}

			condition := metalv1alpha1.AssetCondition{
				Asset:  assetTask.Asset,
				Status: status,
				Type:   "Ready",
			}

			mu.Lock()
			conditions = append(conditions, condition)
			mu.Unlock()
		}

		saveAsset := func(file string) {
			wg.Add(1)

			go func() {
				defer wg.Done()

				if err := save(ctx, assetTask.Asset, file); err != nil {
					setReady(false)

					result = multierror.Append(result, fmt.Errorf("error saving %q: %w", assetTask.Asset.URL, err))
				}

				setReady(true)
				l.Info("saved asset", "url", assetTask.Asset.URL)
			}()
		}

		if _, err := os.Stat(file); os.IsNotExist(err) {
			l.Info("saving asset", "url", assetTask.Asset.URL)
			saveAsset(file)

			continue
		}

		// If we reach here, the file derived from the URL exists, and now we need
		// to update it if the URL has changed.

		l.Info("checking if update required", "file", file)

		ready := false

		for _, condition := range env.Status.Conditions {
			if assetTask.Asset.URL == condition.URL {
				ready = true
			}
		}

		if ready {
			l.Info("update not required", "file", file)
			setReady(true)

			continue
		}

		l.Info("update required", "file", file)

		// At this point the file exists, but the URL for the file has changed. We
		// need to update the file using the new URL.

		l.Info("updating asset", "url", assetTask.Asset.URL)
		saveAsset(file)
	}

	wg.Wait()

	if result.ErrorOrNil() != nil {
		return ctrl.Result{}, result.ErrorOrNil()
	}

	env.Status.Conditions = conditions

	if err := r.Status().Update(ctx, &env); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func save(ctx context.Context, asset metalv1alpha1.Asset, file string) error {
	url := asset.URL

	if url == "" {
		return errors.New("missing URL")
	}

	requestContext, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(requestContext, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		w, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0o666)
		if err != nil {
			return err
		}

		defer w.Close()

		r := resp.Body

		if _, err := io.Copy(w, r); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("failed to download asset: %d", resp.StatusCode)
	}

	return nil
}
