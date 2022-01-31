// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

//nolint:golint,stylecheck
package v1alpha1

import (
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	metalv1alpha2 "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha2"
)

// ConvertTo converts this Environment to the Hub version (v1alpha2).
func (src *Environment) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*metalv1alpha2.Environment)
	if err := Convert_v1alpha1_Environment_To_v1alpha2_Environment(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data from annotations
	restored := &metalv1alpha2.Environment{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	return nil
}

// ConvertFrom converts from the Hub version (v1alpha3) to this version.
func (dst *Environment) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*metalv1alpha2.Environment)
	if err := Convert_v1alpha2_Environment_To_v1alpha1_Environment(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts this MetalMachineTemplateList to the Hub version (v1alpha3).
func (src *EnvironmentList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*metalv1alpha2.EnvironmentList)
	return Convert_v1alpha1_EnvironmentList_To_v1alpha2_EnvironmentList(src, dst, nil)
}

// ConvertFrom converts from the Hub version (v1alpha3) to this version.
func (dst *EnvironmentList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*metalv1alpha2.EnvironmentList)
	return Convert_v1alpha2_EnvironmentList_To_v1alpha1_EnvironmentList(src, dst, nil)
}
