// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// nolint: golint,stylecheck
package v1alpha2

import (
	"log"

	apiconversion "k8s.io/apimachinery/pkg/conversion"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	infrav1alpha3 "github.com/talos-systems/sidero/app/cluster-api-provider-sidero/api/v1alpha3"
)

// ConvertTo converts this MetalCluster to the Hub version (v1alpha3).
func (src *MetalCluster) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1alpha3.MetalCluster)

	if err := Convert_v1alpha2_MetalCluster_To_v1alpha3_MetalCluster(src, dst, nil); err != nil {
		return err
	}

	// Manually convert Spec.APIEndpoints to Spec.ControlPlaneEndpoint.
	log.Printf("length of spec endpoints: %d", len(src.Spec.APIEndpoints))

	if len(src.Spec.APIEndpoints) > 0 {
		endpoint := src.Spec.APIEndpoints[0]
		dst.Spec.ControlPlaneEndpoint.Host = endpoint.Host
		dst.Spec.ControlPlaneEndpoint.Port = int32(endpoint.Port)
	} else if len(src.Status.APIEndpoints) > 0 {
		endpoint := src.Status.APIEndpoints[0]
		dst.Spec.ControlPlaneEndpoint.Host = endpoint.Host
		dst.Spec.ControlPlaneEndpoint.Port = int32(endpoint.Port)
	}

	// Manually restore data.
	restored := &infrav1alpha3.MetalCluster{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	return nil
}

// ConvertFrom converts from the Hub version (v1alpha3) to this version.
func (dst *MetalCluster) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1alpha3.MetalCluster)

	if err := Convert_v1alpha3_MetalCluster_To_v1alpha2_MetalCluster(src, dst, nil); err != nil {
		return err
	}

	// Manually convert Spec.ControlPlaneEndpoint to Status.APIEndpoints.
	if !src.Spec.ControlPlaneEndpoint.IsZero() {
		dst.Status.APIEndpoints = []APIEndpoint{
			{
				Host: src.Spec.ControlPlaneEndpoint.Host,
				Port: int(src.Spec.ControlPlaneEndpoint.Port),
			},
		}
		dst.Spec.APIEndpoints = []APIEndpoint{
			{
				Host: src.Spec.ControlPlaneEndpoint.Host,
				Port: int(src.Spec.ControlPlaneEndpoint.Port),
			},
		}
	}

	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts this MetalClusterList to the Hub version (v1alpha3).
func (src *MetalClusterList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1alpha3.MetalClusterList)
	return Convert_v1alpha2_MetalClusterList_To_v1alpha3_MetalClusterList(src, dst, nil)
}

// ConvertFrom converts from the Hub version (v1alpha3) to this version.
func (dst *MetalClusterList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1alpha3.MetalClusterList)
	return Convert_v1alpha3_MetalClusterList_To_v1alpha2_MetalClusterList(src, dst, nil)
}

// Convert_v1alpha2_MetalClusterStatus_To_v1alpha3_MetalClusterStatus converts MetalCluster.Status from v1alpha2 to v1alpha3.
func Convert_v1alpha2_MetalClusterStatus_To_v1alpha3_MetalClusterStatus(in *MetalClusterStatus, out *infrav1alpha3.MetalClusterStatus, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha2_MetalClusterStatus_To_v1alpha3_MetalClusterStatus(in, out, s); err != nil {
		return err
	}

	return nil
}

// Convert_v1alpha2_MetalClusterSpec_To_v1alpha3_MetalClusterSpec.
func Convert_v1alpha2_MetalClusterSpec_To_v1alpha3_MetalClusterSpec(in *MetalClusterSpec, out *infrav1alpha3.MetalClusterSpec, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha2_MetalClusterSpec_To_v1alpha3_MetalClusterSpec(in, out, s); err != nil {
		return err
	}

	return nil
}

// Convert_v1alpha3_MetalClusterSpec_To_v1alpha2_MetalClusterSpec converts from the Hub version (v1alpha3) of the MetalClusterSpec to this version.
func Convert_v1alpha3_MetalClusterSpec_To_v1alpha2_MetalClusterSpec(in *infrav1alpha3.MetalClusterSpec, out *MetalClusterSpec, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha3_MetalClusterSpec_To_v1alpha2_MetalClusterSpec(in, out, s); err != nil {
		return err
	}

	return nil
}

// Convert_v1alpha3_MetalClusterStatus_To_v1alpha2_MetalClusterStatus.
func Convert_v1alpha3_MetalClusterStatus_To_v1alpha2_MetalClusterStatus(in *infrav1alpha3.MetalClusterStatus, out *MetalClusterStatus, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha3_MetalClusterStatus_To_v1alpha2_MetalClusterStatus(in, out, s); err != nil {
		return err
	}

	return nil
}
