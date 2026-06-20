// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha2

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func (r *ServerClass) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(r).
		Complete()
}

func (r *ServerClassList) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:verbs=create;update;delete,path=/validate-metal-sidero-dev-v1alpha2-serverclass,mutating=false,failurePolicy=fail,groups=metal.sidero.dev,resources=serverclasses,versions=v1alpha2,name=vserverclasses.metal.sidero.dev,sideEffects=None,admissionReviewVersions=v1

var _ webhook.CustomValidator = &ServerClass{}

func (r *ServerClass) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	r = obj.(*ServerClass)

	return nil, r.validate()
}

func (r *ServerClass) ValidateUpdate(ctx context.Context, oldObj runtime.Object, newObj runtime.Object) (admission.Warnings, error) {
	r = newObj.(*ServerClass)

	return nil, r.validate()
}

func (r *ServerClass) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func (r *ServerClass) validate() error {
	allErrs := validateInstallDisk(field.NewPath("spec").Child("installDisk"), r.Spec.InstallDisk, r.Spec.WipeOnlyInstallDisk)
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(schema.GroupKind{Group: GroupVersion.Group, Kind: "ServerClass"}, r.Name, allErrs)
}
