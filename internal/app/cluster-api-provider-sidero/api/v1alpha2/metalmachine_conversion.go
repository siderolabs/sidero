// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// nolint: golint,stylecheck
package v1alpha2

import (
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	infrav1alpha3 "github.com/talos-systems/sidero/internal/app/cluster-api-provider-sidero/api/v1alpha3"
)

// ConvertTo converts this MetalMachine to the Hub version (v1alpha3).
func (src *MetalMachine) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1alpha3.MetalMachine)

	if err := Convert_v1alpha2_MetalMachine_To_v1alpha3_MetalMachine(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data from annotations
	restored := &infrav1alpha3.MetalMachine{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	return nil
}

// ConvertFrom converts from the Hub version (v1alpha3) to this version.
func (dst *MetalMachine) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1alpha3.MetalMachine)
	if err := Convert_v1alpha3_MetalMachine_To_v1alpha2_MetalMachine(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts this MetalMachineList to the Hub version (v1alpha3).
func (src *MetalMachineList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1alpha3.MetalMachineList)
	return Convert_v1alpha2_MetalMachineList_To_v1alpha3_MetalMachineList(src, dst, nil)
}

// ConvertFrom converts from the Hub version (v1alpha3) to this version.
func (dst *MetalMachineList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1alpha3.MetalMachineList)
	return Convert_v1alpha3_MetalMachineList_To_v1alpha2_MetalMachineList(src, dst, nil)
}

func Convert_v1alpha2_MetalMachineSpec_To_v1alpha3_MetalMachineSpec(in *MetalMachineSpec, out *infrav1alpha3.MetalMachineSpec, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha2_MetalMachineSpec_To_v1alpha3_MetalMachineSpec(in, out, s); err != nil {
		return err
	}

	return nil
}

// Convert_v1alpha3_MetalMachineSpec_To_v1alpha2_MetalMachineSpec converts from the Hub version (v1alpha3) of the MetalMachineSpec to this version.
func Convert_v1alpha3_MetalMachineSpec_To_v1alpha2_MetalMachineSpec(in *infrav1alpha3.MetalMachineSpec, out *MetalMachineSpec, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha3_MetalMachineSpec_To_v1alpha2_MetalMachineSpec(in, out, s); err != nil {
		return err
	}

	return nil
}

// Convert_v1alpha2_MetalMachineStatus_To_v1alpha3_MetalMachineStatus converts this MetalMachineStatus to the Hub version (v1alpha3).
func Convert_v1alpha2_MetalMachineStatus_To_v1alpha3_MetalMachineStatus(in *MetalMachineStatus, out *infrav1alpha3.MetalMachineStatus, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha2_MetalMachineStatus_To_v1alpha3_MetalMachineStatus(in, out, s); err != nil {
		return err
	}

	// Manually convert the Error fields to the Failure fields
	out.FailureMessage = in.ErrorMessage
	out.FailureReason = in.ErrorReason

	return nil
}

// Convert_v1alpha3_MetalMachineStatus_To_v1alpha2_MetalMachineStatus converts from the Hub version (v1alpha3) of the MetalMachineStatus to this version.
func Convert_v1alpha3_MetalMachineStatus_To_v1alpha2_MetalMachineStatus(in *infrav1alpha3.MetalMachineStatus, out *MetalMachineStatus, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha3_MetalMachineStatus_To_v1alpha2_MetalMachineStatus(in, out, s); err != nil {
		return err
	}

	// Manually convert the Failure fields to the Error fields
	out.ErrorMessage = in.FailureMessage
	out.ErrorReason = in.FailureReason

	return nil
}
