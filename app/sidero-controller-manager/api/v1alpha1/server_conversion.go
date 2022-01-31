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

// ConvertTo converts this Server to the Hub version (v1alpha2).
func (src *Server) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*metalv1alpha2.Server)
	if err := Convert_v1alpha1_Server_To_v1alpha2_Server(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data from annotations
	restored := &metalv1alpha2.Server{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	return nil
}

// ConvertFrom converts from the Hub version (v1alpha3) to this version.
func (dst *Server) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*metalv1alpha2.Server)
	if err := Convert_v1alpha2_Server_To_v1alpha1_Server(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts this MetalMachineTemplateList to the Hub version (v1alpha3).
func (src *ServerList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*metalv1alpha2.ServerList)
	return Convert_v1alpha1_ServerList_To_v1alpha2_ServerList(src, dst, nil)
}

// ConvertFrom converts from the Hub version (v1alpha3) to this version.
func (dst *ServerList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*metalv1alpha2.ServerList)
	return Convert_v1alpha2_ServerList_To_v1alpha1_ServerList(src, dst, nil)
}

// Convert_v1alpha1_ServerSpec_To_v1alpha2_ServerSpec converts to the Hub version (v1alpha2).
func Convert_v1alpha1_ServerSpec_To_v1alpha2_ServerSpec(in *ServerSpec, out *metalv1alpha2.ServerSpec, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha1_ServerSpec_To_v1alpha2_ServerSpec(in, out, s); err != nil {
		return err
	}

	// Manually convert SystemInformation to Hardware.
	if in.SystemInformation != nil {
		if out.Hardware == nil {
			out.Hardware = &metalv1alpha2.HardwareInformation{}
		}
		out.Hardware.System = &metalv1alpha2.SystemInformation{
			Manufacturer: in.SystemInformation.Manufacturer,
			ProductName:  in.SystemInformation.ProductName,
			Version:      in.SystemInformation.Version,
			SerialNumber: in.SystemInformation.SerialNumber,
			SKUNumber:    in.SystemInformation.SKUNumber,
			Family:       in.SystemInformation.Family,
		}
	}

	// Manually convert CPU to Hardware.
	if in.CPU != nil {
		if out.Hardware == nil {
			out.Hardware = &metalv1alpha2.HardwareInformation{}
		}
		out.Hardware.Compute = &metalv1alpha2.ComputeInformation{
			Processors: []*metalv1alpha2.Processor{
				{
					Manufacturer: in.CPU.Manufacturer,
					ProductName:  in.CPU.Version,
				},
			},
		}
	}

	return nil
}

// Convert_v1alpha2_ServerSpec_To_v1alpha1_ServerSpec converts from the Hub version (v1alpha2).
func Convert_v1alpha2_ServerSpec_To_v1alpha1_ServerSpec(in *metalv1alpha2.ServerSpec, out *ServerSpec, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha2_ServerSpec_To_v1alpha1_ServerSpec(in, out, s); err != nil {
		return err
	}

	// Manually convert Hardware to SystemInformation.
	if in.Hardware != nil && in.Hardware.System != nil {
		out.SystemInformation = &SystemInformation{
			Manufacturer: in.Hardware.System.Manufacturer,
			ProductName:  in.Hardware.System.ProductName,
			Version:      in.Hardware.System.Version,
			SerialNumber: in.Hardware.System.SerialNumber,
			SKUNumber:    in.Hardware.System.SKUNumber,
			Family:       in.Hardware.System.Family,
		}
	}

	// Manually convert Hardware to CPU.
	if in.Hardware != nil && in.Hardware.Compute != nil && len(in.Hardware.Compute.Processors) > 0 {
		cpu := in.Hardware.Compute.Processors[0]
		out.CPU = &CPUInformation{
			Manufacturer: cpu.Manufacturer,
			Version:      cpu.ProductName,
		}
	}

	return nil
}

func Convert_v1alpha2_SystemInformation_To_v1alpha1_SystemInformation(in *metalv1alpha2.SystemInformation, out *SystemInformation, s apiconversion.Scope) error {
	return autoConvert_v1alpha2_SystemInformation_To_v1alpha1_SystemInformation(in, out, s)
}
