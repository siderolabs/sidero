// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

//nolint:golint,stylecheck
package v1alpha1

import (
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	metalv1alpha2 "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha2"
)

// ConvertTo converts this ServerClass to the Hub version (v1alpha2).
func (src *ServerClass) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*metalv1alpha2.ServerClass)
	if err := Convert_v1alpha1_ServerClass_To_v1alpha2_ServerClass(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data from annotations
	restored := &metalv1alpha2.ServerClass{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	return nil
}

// ConvertFrom converts from the Hub version (v1alpha3) to this version.
func (dst *ServerClass) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*metalv1alpha2.ServerClass)
	if err := Convert_v1alpha2_ServerClass_To_v1alpha1_ServerClass(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts this MetalMachineTemplateList to the Hub version (v1alpha3).
func (src *ServerClassList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*metalv1alpha2.ServerClassList)
	return Convert_v1alpha1_ServerClassList_To_v1alpha2_ServerClassList(src, dst, nil)
}

// ConvertFrom converts from the Hub version (v1alpha3) to this version.
func (dst *ServerClassList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*metalv1alpha2.ServerClassList)
	return Convert_v1alpha2_ServerClassList_To_v1alpha1_ServerClassList(src, dst, nil)
}

// Convert_v1alpha1_Qualifiers_To_v1alpha2_Qualifiers converts to the Hub version (v1alpha2).
func Convert_v1alpha1_Qualifiers_To_v1alpha2_Qualifiers(in *Qualifiers, out *metalv1alpha2.Qualifiers, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha1_Qualifiers_To_v1alpha2_Qualifiers(in, out, s); err != nil {
		return err
	}

	// Manually convert SystemInformation to Hardware.
	for _, v := range in.SystemInformation {
		out.Hardware = append(out.Hardware, metalv1alpha2.HardwareInformation{
			System: &metalv1alpha2.SystemInformation{
				Manufacturer: v.Manufacturer,
				ProductName:  v.ProductName,
				Version:      v.Version,
				SerialNumber: v.SerialNumber,
				SKUNumber:    v.SKUNumber,
				Family:       v.Family,
			},
		})
	}

	// Manually convert CPU to Hardware.
	for _, v := range in.CPU {
		out.Hardware = append(out.Hardware, metalv1alpha2.HardwareInformation{
			Compute: &metalv1alpha2.ComputeInformation{
				Processors: []*metalv1alpha2.Processor{
					{
						Manufacturer: v.Manufacturer,
						ProductName:  v.Version,
					},
				},
			},
		})
	}

	return nil
}

// Convert_v1alpha2_Qualifiers_To_v1alpha1_Qualifiers converts from the Hub version (v1alpha2).
func Convert_v1alpha2_Qualifiers_To_v1alpha1_Qualifiers(in *metalv1alpha2.Qualifiers, out *Qualifiers, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha2_Qualifiers_To_v1alpha1_Qualifiers(in, out, s); err != nil {
		return err
	}

	// Manually convert Hardware to SystemInformation or CPU.
	for _, v := range in.Hardware {
		if v.System != nil {
			out.SystemInformation = append(out.SystemInformation, SystemInformation{
				Manufacturer: v.System.Manufacturer,
				ProductName:  v.System.ProductName,
				Version:      v.System.Version,
				SerialNumber: v.System.SerialNumber,
				SKUNumber:    v.System.SKUNumber,
				Family:       v.System.Family,
			})
		}
		if v.Compute != nil && len(v.Compute.Processors) > 0 {
			cpu := v.Compute.Processors[0]
			out.CPU = append(out.CPU, CPUInformation{
				Manufacturer: cpu.Manufacturer,
				Version:      cpu.ProductName,
			})
		}
	}

	return nil
}
