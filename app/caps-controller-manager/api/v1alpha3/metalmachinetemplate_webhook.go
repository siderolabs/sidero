// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha3

import (
	"github.com/google/go-cmp/cmp"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	runtime "k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func (r *MetalMachineTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:verbs=create;update;delete,path=/validate-infrastructure-cluster-x-k8s-io-v1alpha3-metalmachinetemplate,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=metalmachinetemplates,versions=v1alpha3,name=vmetalmachinetemplates.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1

var _ webhook.Validator = &MetalMachineTemplate{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *MetalMachineTemplate) ValidateCreate() (admission.Warnings, error) {
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *MetalMachineTemplate) ValidateUpdate(oldRaw runtime.Object) (admission.Warnings, error) {
	old := oldRaw.(*MetalMachineTemplate)

	if !cmp.Equal(r.Spec, old.Spec) {
		return nil, apierrors.NewBadRequest("MetalMachineTemplate.Spec is immutable")
	}

	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *MetalMachineTemplate) ValidateDelete() (admission.Warnings, error) {
	return nil, nil
}
