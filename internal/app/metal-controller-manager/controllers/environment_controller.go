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
	"github.com/hashicorp/go-multierror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	metalv1alpha1 "github.com/talos-systems/sidero/internal/app/metal-controller-manager/api/v1alpha1"
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
		assets     = []metalv1alpha1.Asset{env.Spec.Kernel.Asset, env.Spec.Initrd.Asset}
		conditions = []metalv1alpha1.AssetCondition{}
		wg         sync.WaitGroup
		mu         sync.Mutex
		result     *multierror.Error
	)

	for _, asset := range assets {
		asset := asset

		file := filepath.Join(envs, filepath.Base(asset.URL))

		if _, err := os.Stat(file); os.IsNotExist(err) {
			wg.Add(1)

			go func() {
				defer wg.Done()

				l.Info("saving asset", "url", asset.URL)

				if err := save(asset, file); err != nil {
					condition := metalv1alpha1.AssetCondition{
						Asset:  asset,
						Status: "False",
						Type:   "Ready",
					}

					mu.Lock()
					conditions = append(conditions, condition)
					mu.Unlock()

					result = multierror.Append(result, fmt.Errorf("error saving %q: %w", asset.URL, err))
				}

				l.Info("saved asset", "url", asset.URL)

				condition := metalv1alpha1.AssetCondition{
					Asset:  asset,
					Status: "True",
					Type:   "Ready",
				}

				mu.Lock()
				conditions = append(conditions, condition)
				mu.Unlock()
			}()

			continue
		}

		// If we reach here, the file derived from the URL exists, and now we need
		// to update it if the URL has changed.

		l.Info("checking if update required", "file", file)

		ready := false

		for _, condition := range env.Status.Conditions {
			if asset.URL == condition.URL {
				ready = true
			}
		}

		if ready {
			l.Info("update not required", "file", file)

			condition := metalv1alpha1.AssetCondition{
				Asset:  asset,
				Status: "True",
				Type:   "Ready",
			}

			conditions = append(conditions, condition)

			continue
		}

		l.Info("update required", "file", file)

		// At this point the file exists, but the URL for the file has changed. We
		// need to update the file using the new URL.

		wg.Add(1)

		go func() {
			defer wg.Done()

			l.Info("updating asset", "url", asset.URL)

			if err := save(asset, file); err != nil {
				condition := metalv1alpha1.AssetCondition{
					Asset:  asset,
					Status: "False",
					Type:   "Ready",
				}

				mu.Lock()
				conditions = append(conditions, condition)
				mu.Unlock()

				result = multierror.Append(result, fmt.Errorf("error updating %q: %w", asset.URL, err))
			}

			l.Info("updated asset", "url", asset.URL)

			condition := metalv1alpha1.AssetCondition{
				Asset:  asset,
				Status: "True",
				Type:   "Ready",
			}

			mu.Lock()
			conditions = append(conditions, condition)
			mu.Unlock()
		}()
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

func save(asset metalv1alpha1.Asset, file string) error {
	url := asset.URL

	if url == "" {
		return errors.New("missing URL")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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
