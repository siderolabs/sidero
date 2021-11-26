// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"inet.af/netaddr"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/talos-systems/siderolink/pkg/logreceiver"

	sidero "github.com/talos-systems/sidero/app/caps-controller-manager/api/v1alpha3"
)

type sourceAnnotation struct {
	ServerUUID       string
	MetalMachineName string
	MachineName      string
	ClusterName      string
}

var sourceMap sync.Map

func fetchSourceAnnotation(ctx context.Context, metalClient runtimeclient.Client, srcAddr netaddr.IP) (sourceAnnotation, error) {
	var (
		annotation     sourceAnnotation
		serverbindings sidero.ServerBindingList
	)

	if err := metalClient.List(ctx, &serverbindings); err != nil {
		return annotation, fmt.Errorf("error getting server bindings: %w", err)
	}

	srcAddress := srcAddr.String()

	var serverBinding *sidero.ServerBinding

	for _, item := range serverbindings.Items {
		item := item

		if strings.HasPrefix(item.Spec.SideroLink.NodeAddress, srcAddress) {
			serverBinding = &item

			break
		}
	}

	if serverBinding == nil {
		// no matching server binding, leave things as is
		return annotation, nil
	}

	annotation.ServerUUID = serverBinding.Name
	annotation.MetalMachineName = fmt.Sprintf("%s/%s", serverBinding.Spec.MetalMachineRef.Namespace, serverBinding.Spec.MetalMachineRef.Name)
	annotation.ClusterName = serverBinding.Labels[clusterv1.ClusterLabelName]

	var metalMachine sidero.MetalMachine

	if err := metalClient.Get(ctx,
		types.NamespacedName{
			Namespace: serverBinding.Spec.MetalMachineRef.Namespace,
			Name:      serverBinding.Spec.MetalMachineRef.Name,
		},
		&metalMachine); err != nil {
		return annotation, fmt.Errorf("error getting metal machine: %w", err)
	}

	for _, ref := range metalMachine.OwnerReferences {
		gv, err := schema.ParseGroupVersion(ref.APIVersion)
		if err != nil {
			continue
		}

		if ref.Kind == "Machine" && gv.Group == clusterv1.GroupVersion.Group {
			annotation.MachineName = fmt.Sprintf("%s/%s", metalMachine.Namespace, ref.Name)

			break
		}
	}

	return annotation, nil
}

func logHandler(metalClient runtimeclient.Client, logger *zap.Logger) logreceiver.Handler {
	return func(srcAddr netaddr.IP, msg map[string]interface{}) {
		var annotation sourceAnnotation

		v, ok := sourceMap.Load(srcAddr)
		if !ok {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			var err error

			annotation, err = fetchSourceAnnotation(ctx, metalClient, srcAddr)
			if err != nil {
				logger.Error("error fetching server information", zap.Error(err), zap.Stringer("source_addr", srcAddr))
			}

			sourceMap.Store(srcAddr, annotation)
		} else {
			annotation = v.(sourceAnnotation)
		}

		if annotation.ServerUUID != "" {
			msg["server_uuid"] = annotation.ServerUUID
		}

		if annotation.ClusterName != "" {
			msg["cluster"] = annotation.ClusterName
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
