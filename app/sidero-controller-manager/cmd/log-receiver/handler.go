// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"encoding/json"
	"net/netip"
	"os"

	"go.uber.org/zap"

	"github.com/siderolabs/sidero/app/sidero-controller-manager/internal/siderolink"
	"github.com/siderolabs/siderolink/pkg/logreceiver"
)

func logHandler(logger *zap.Logger, annotator *siderolink.Annotator) logreceiver.Handler {
	return func(srcAddr netip.Addr, msg map[string]interface{}) {
		annotation, _ := annotator.Get(srcAddr.String())

		if annotation.ServerUUID != "" {
			msg["server_uuid"] = annotation.ServerUUID
		}

		if annotation.ClusterName != "" {
			msg["cluster"] = annotation.ClusterName
		}

		if annotation.Namespace != "" {
			msg["namespace"] = annotation.Namespace
		}

		if annotation.MetalMachineName != "" {
			msg["metal_machine"] = annotation.MetalMachineName
		}

		if annotation.MachineName != "" {
			msg["machine"] = annotation.MachineName
		}

		if err := json.NewEncoder(os.Stdout).Encode(msg); err != nil {
			logger.Error("error printing log message", zap.Error(err))
		}
	}
}
