// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha1

import (
	"fmt"
	"sort"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	siderotypes "github.com/talos-systems/sidero/app/sidero-controller-manager/pkg/types"
)

var operations = map[string]struct{}{
	"add":     {},
	"remove":  {},
	"replace": {},
	"copy":    {},
	"move":    {},
	"test":    {},
}

var operationKinds = []string{}

func (r *Server) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:verbs=create;update;delete,path=/validate-metal-sidero-dev-v1alpha1-server,mutating=false,failurePolicy=fail,groups=metal.sidero.dev,resources=servers,versions=v1alpha1,name=vservers.metal.sidero.dev,sideEffects=None,admissionReviewVersions=v1

var _ webhook.Validator = &Server{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *Server) ValidateCreate() error {
	return r.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *Server) ValidateUpdate(old runtime.Object) error {
	return r.validate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *Server) ValidateDelete() error {
	return nil
}

func (r *Server) validate() error {
	var allErrs field.ErrorList

	allErrs = append(allErrs, r.validateBootFromDisk()...)
	allErrs = append(allErrs, r.validatePXEMode()...)
	allErrs = append(allErrs, r.validateConfigPatches()...)

	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: GroupVersion.Group, Kind: "Server"},
		r.Name, allErrs)
}

func (r *Server) validateBootFromDisk() (allErrs field.ErrorList) {
	validValues := []siderotypes.BootFromDisk{
		"",
		siderotypes.BootIPXEExit,
		siderotypes.Boot404,
		siderotypes.BootSANDisk,
	}

	var valid bool

	for _, v := range validValues {
		if r.Spec.BootFromDiskMethod == v {
			valid = true

			break
		}
	}

	if !valid {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec").Child("bootFromDiskMethod"), r.Spec.BootFromDiskMethod,
				fmt.Sprintf("valid values are: %q", validValues),
			),
		)
	}

	return allErrs
}

func (r *Server) validatePXEMode() (allErrs field.ErrorList) {
	validValues := []siderotypes.PXEMode{
		"",
		siderotypes.PXEModeBIOS,
		siderotypes.PXEModeUEFI,
	}

	var valid bool

	for _, v := range validValues {
		if r.Spec.PXEMode == v {
			valid = true

			break
		}
	}

	if !valid {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec").Child("pxeMode"), r.Spec.BootFromDiskMethod,
				fmt.Sprintf("valid values are: %q", validValues),
			),
		)
	}

	return allErrs
}

func (r *Server) validateConfigPatches() (allErrs field.ErrorList) {
	for index, patch := range r.Spec.ConfigPatches {
		if _, ok := operations[patch.Op]; !ok {
			allErrs = append(allErrs,
				field.Invalid(field.NewPath("spec").Child("configPatches").Child(fmt.Sprintf("%d", index)).Child("op"), patch.Op,
					fmt.Sprintf("valid values are: %q", operationKinds),
				),
			)
		}
	}

	return allErrs
}

func init() {
	operationKinds = make([]string, 0, len(operations))

	for key := range operations {
		operationKinds = append(operationKinds, key)
	}

	sort.Strings(operationKinds)
}
