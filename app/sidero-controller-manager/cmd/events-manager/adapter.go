// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"inet.af/netaddr"
	"k8s.io/apimachinery/pkg/types"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	sidero "github.com/talos-systems/sidero/app/caps-controller-manager/api/v1alpha3"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/siderolink"

	"github.com/talos-systems/siderolink/pkg/events"

	"github.com/talos-systems/talos/pkg/machinery/api/machine"
)

// Adapter implents gRPC API.
type Adapter struct {
	Sink *events.Sink

	logger      *zap.Logger
	annotator   *siderolink.Annotator
	metalClient runtimeclient.Client
}

// NewAdapter initializes new server.
func NewAdapter(metalClient runtimeclient.Client, annotator *siderolink.Annotator, logger *zap.Logger) *Adapter {
	return &Adapter{
		logger:      logger,
		annotator:   annotator,
		metalClient: metalClient,
	}
}

// HandleEvent implements events.Adapter.
func (a *Adapter) HandleEvent(ctx context.Context, event events.Event) error {
	logger := a.logger.With(
		zap.String("node", event.Node),
		zap.String("id", event.ID),
		zap.String("type", event.TypeURL),
	)

	fields := []zap.Field{}
	message := "incoming event"

	var err error

	ipPort, err := netaddr.ParseIPPort(event.Node)
	if err != nil {
		return err
	}

	ip := ipPort.IP().String()

	annotation, _ := a.annotator.Get(ip)

	if annotation.ServerUUID != "" {
		fields = append(fields, zap.String("server_uuid", annotation.ServerUUID))
	}

	if annotation.ClusterName != "" {
		fields = append(fields, zap.String("cluster", annotation.ClusterName))
	}

	if annotation.Namespace != "" {
		fields = append(fields, zap.String("namespace", annotation.Namespace))
	}

	if annotation.MetalMachineName != "" {
		fields = append(fields, zap.String("metal_machine", annotation.MetalMachineName))
	}

	if annotation.MachineName != "" {
		fields = append(fields, zap.String("machine", annotation.MachineName))
	}

	switch event := event.Payload.(type) {
	case *machine.AddressEvent:
		fields = append(fields, zap.String("hostname", event.GetHostname()), zap.String("addresses", strings.Join(event.GetAddresses(), ",")))

		// filter out SideroLink address from the list
		addresses := event.Addresses

		n := 0

		for _, addr := range addresses {
			if addr != ip {
				addresses[n] = addr
				n++
			}
		}

		addresses = addresses[:n]

		err = a.patchServerBinding(ctx, ip, func(serverbinding *sidero.ServerBinding) {
			serverbinding.Spec.Addresses = addresses
			serverbinding.Spec.Hostname = event.Hostname
		})

		if err != nil {
			return err
		}
	case *machine.ConfigValidationErrorEvent:
		fields = append(fields, zap.Error(fmt.Errorf(event.GetError())))

		if err = a.handleConfigValidationFailedEvent(ctx, ip, event); err != nil {
			return err
		}
	case *machine.ConfigLoadErrorEvent:
		fields = append(fields, zap.Error(fmt.Errorf(event.GetError())))

		if err = a.handleConfigLoadFailedEvent(ctx, ip, event); err != nil {
			return err
		}
	case *machine.PhaseEvent:
		fields = append(fields, zap.String("phase", event.GetPhase()), zap.String("action", event.GetAction().String()))

		if err = a.handlePhaseEvent(ctx, ip, event); err != nil {
			return err
		}
	case *machine.TaskEvent:
		fields = append(fields, zap.String("task", event.GetTask()), zap.String("action", event.GetAction().String()))
	case *machine.ServiceStateEvent:
		message = "service " + event.GetMessage()
		fields = append(fields, zap.String("service", event.GetService()), zap.String("action", event.GetAction().String()))
	case *machine.SequenceEvent:
		fields = append(fields, zap.String("sequence", event.GetSequence()), zap.String("action", event.GetAction().String()))

		if err = a.handleSequenceEvent(ctx, ip, event); err != nil {
			return err
		}

		if event.GetError() != nil {
			err = fmt.Errorf(event.GetError().GetMessage())
		}

		if event.GetSequence() == "install" &&
			event.GetAction() == machine.SequenceEvent_STOP {
			if event.GetError() != nil {
				message = "failed to install Talos"
				break
			}

			message = "successfully installed Talos"
		}
	}

	if err != nil {
		fields = append(fields, zap.Error(err))

		logger.Error(message, fields...)

		return nil
	}

	logger.Info(message, fields...)

	return nil
}

func (a *Adapter) handleSequenceEvent(ctx context.Context, ip string, event *machine.SequenceEvent) error {
	if event.GetSequence() != "install" {
		return nil
	}

	var callback func(*sidero.ServerBinding)

	if event.GetAction() == machine.SequenceEvent_STOP {
		if event.GetError() != nil {
			callback = func(serverbinding *sidero.ServerBinding) {
				conditions.MarkFalse(serverbinding, sidero.TalosInstalledCondition, sidero.TalosInstallationFailedReason, clusterv1.ConditionSeverityError, event.GetError().GetMessage())
			}
		} else {
			callback = func(serverbinding *sidero.ServerBinding) {
				conditions.MarkTrue(serverbinding, sidero.TalosInstalledCondition)
				conditions.MarkTrue(serverbinding, sidero.TalosConfigValidatedCondition)
				conditions.MarkTrue(serverbinding, sidero.TalosConfigLoadedCondition)
			}
		}
	} else if event.GetAction() == machine.SequenceEvent_START {
		callback = func(serverbinding *sidero.ServerBinding) {
			conditions.MarkFalse(serverbinding, sidero.TalosInstalledCondition, sidero.TalosInstallationInProgressReason, clusterv1.ConditionSeverityInfo, "")
			conditions.MarkFalse(serverbinding, sidero.TalosConfigValidatedCondition, sidero.TalosInstallationInProgressReason, clusterv1.ConditionSeverityInfo, "")
			conditions.MarkFalse(serverbinding, sidero.TalosConfigLoadedCondition, sidero.TalosInstallationInProgressReason, clusterv1.ConditionSeverityInfo, "")
		}
	}

	if callback == nil {
		return nil
	}

	return a.patchServerBinding(ctx, ip, callback)
}

func (a *Adapter) handleConfigLoadFailedEvent(ctx context.Context, ip string, event *machine.ConfigLoadErrorEvent) error {
	return a.patchServerBinding(ctx, ip, func(serverbinding *sidero.ServerBinding) {
		conditions.MarkFalse(serverbinding, sidero.TalosConfigLoadedCondition, sidero.TalosConfigLoadFailedReason, clusterv1.ConditionSeverityError, event.GetError())
	})
}

func (a *Adapter) handleConfigValidationFailedEvent(ctx context.Context, ip string, event *machine.ConfigValidationErrorEvent) error {
	return a.patchServerBinding(ctx, ip, func(serverbinding *sidero.ServerBinding) {
		conditions.MarkFalse(serverbinding, sidero.TalosConfigValidatedCondition, sidero.TalosConfigValidationFailedReason, clusterv1.ConditionSeverityError, event.GetError())
	})
}

func (a *Adapter) handlePhaseEvent(ctx context.Context, ip string, event *machine.PhaseEvent) (err error) {
	if event.GetAction() != machine.PhaseEvent_STOP {
		return nil
	}

	switch event.GetPhase() {
	case "validateConfig":
		if err = a.patchServerBinding(ctx, ip, func(serverbinding *sidero.ServerBinding) {
			if !conditions.Has(serverbinding, sidero.TalosConfigValidatedCondition) {
				conditions.MarkTrue(serverbinding, sidero.TalosConfigValidatedCondition)
			}
		}); err != nil {
			return err
		}
	case "loadConfig":
		if err = a.patchServerBinding(ctx, ip, func(serverbinding *sidero.ServerBinding) {
			if !conditions.Has(serverbinding, sidero.TalosConfigLoadedCondition) {
				conditions.MarkTrue(serverbinding, sidero.TalosConfigLoadedCondition)
			}
		}); err != nil {
			return err
		}
	}

	return nil
}

func (a *Adapter) patchServerBinding(ctx context.Context, ip string, callback func(serverbinding *sidero.ServerBinding)) error {
	annotation, exists := a.annotator.Get(ip)
	if !exists {
		return fmt.Errorf("failed to find ServerBindings for ip %s", ip)
	}

	var serverbinding sidero.ServerBinding
	if err := a.metalClient.Get(ctx, types.NamespacedName{Name: annotation.ServerUUID}, &serverbinding); err != nil {
		return err
	}

	patchHelper, err := patch.NewHelper(&serverbinding, a.metalClient)
	if err != nil {
		return err
	}

	callback(&serverbinding)

	return patchHelper.Patch(ctx, &serverbinding)
}
