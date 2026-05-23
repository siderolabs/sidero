// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha2

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

func validateInstallDisk(path *field.Path, installDisk *InstallDisk, wipeOnlyInstallDisk bool) (allErrs field.ErrorList) {
	if installDisk == nil {
		if wipeOnlyInstallDisk {
			allErrs = append(allErrs, field.Required(path, "must be set when spec.wipeOnlyInstallDisk is true"))
		}

		return allErrs
	}

	hasDevice := installDisk.DeviceName != ""
	hasSelector := installDisk.Selector != nil

	if hasDevice == hasSelector {
		allErrs = append(allErrs, field.Invalid(path, installDisk, "exactly one of deviceName or selector must be set"))
	}

	if hasSelector {
		sel := installDisk.Selector
		if sel.Name == "" && sel.Model == "" && sel.Serial == "" && sel.UUID == "" && sel.WWID == "" && sel.Type == "" && sel.Size == 0 {
			allErrs = append(allErrs, field.Required(path.Child("selector"), "must not be empty"))
		}
		if sel.Type != "" {
			switch sel.Type {
			case "ssd", "hdd", "nvme", "sd":
			default:
				allErrs = append(allErrs, field.Invalid(path.Child("selector").Child("type"), sel.Type, fmt.Sprintf("valid values are: %q", []string{"ssd", "hdd", "nvme", "sd"})))
			}
		}
	}

	return allErrs
}
