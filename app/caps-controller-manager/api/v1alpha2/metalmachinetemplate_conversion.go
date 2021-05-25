// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

//nolint:golint,stylecheck
package v1alpha2

import (
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	infrav1alpha3 "github.com/talos-systems/sidero/app/caps-controller-manager/api/v1alpha3"
)

// ConvertTo converts this MetalMachineTemplate to the Hub version (v1alpha3).
func (src *MetalMachineTemplate) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1alpha3.MetalMachineTemplate)
	if err := Convert_v1alpha2_MetalMachineTemplate_To_v1alpha3_MetalMachineTemplate(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data from annotations
	restored := &infrav1alpha3.MetalMachineTemplate{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	return nil
}

// ConvertFrom converts from the Hub version (v1alpha3) to this version.
func (dst *MetalMachineTemplate) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1alpha3.MetalMachineTemplate)
	if err := Convert_v1alpha3_MetalMachineTemplate_To_v1alpha2_MetalMachineTemplate(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts this MetalMachineTemplateList to the Hub version (v1alpha3).
func (src *MetalMachineTemplateList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1alpha3.MetalMachineTemplateList)
	return Convert_v1alpha2_MetalMachineTemplateList_To_v1alpha3_MetalMachineTemplateList(src, dst, nil)
}

// ConvertFrom converts from the Hub version (v1alpha3) to this version.
func (dst *MetalMachineTemplateList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1alpha3.MetalMachineTemplateList)
	return Convert_v1alpha3_MetalMachineTemplateList_To_v1alpha2_MetalMachineTemplateList(src, dst, nil)
}
