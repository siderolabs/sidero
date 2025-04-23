// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"strings"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/siderolabs/gen/xslices"

	sidero "github.com/siderolabs/sidero/app/caps-controller-manager/api/v1alpha3"
	"github.com/siderolabs/sidero/app/sidero-controller-manager/internal/siderolink"
	"github.com/siderolabs/siderolink/pkg/events"

	"github.com/siderolabs/talos/pkg/machinery/api/common"
	"github.com/siderolabs/talos/pkg/machinery/api/machine"
)

// Adapter implents gRPC API.
type Adapter struct {
	Sink *events.Sink

	logger                *zap.Logger
	annotator             *siderolink.Annotator
	metalClient           runtimeclient.Client
	negativeAddressFilter []netip.Prefix
}

// NewAdapter initializes new server.
func NewAdapter(metalClient runtimeclient.Client, annotator *siderolink.Annotator, logger *zap.Logger, negativeAddressFilter []netip.Prefix) *Adapter {
	return &Adapter{
		logger:                logger,
		annotator:             annotator,
		metalClient:           metalClient,
		negativeAddressFilter: negativeAddressFilter,
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

	ipPort, err := netip.ParseAddrPort(event.Node)
	if err != nil {
		return err
	}

	ip := ipPort.Addr()

	annotation, _ := a.annotator.Get(ip.String())

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
		fields = append(fields, zap.String("hostname", event.GetHostname()), zap.Strings("addresses", event.GetAddresses()))

		// filter out SideroLink address and other negative matches from the list
		addresses := xslices.Filter(event.Addresses,
			func(addrStr string) bool {
				addr, err := netip.ParseAddr(addrStr)
				if err != nil {
					// invalid address
					return false
				}

				if addr == ip {
					// SideroLink address
					return false
				}

				for _, prefix := range a.negativeAddressFilter {
					if prefix.Contains(addr) {
						return false
					}
				}

				return true
			})

		err = a.patchServerBinding(ctx, ip.String(), func(serverbinding *sidero.ServerBinding) {
			serverbinding.Spec.Addresses = addresses
			serverbinding.Spec.Hostname = event.Hostname
		})
		if err != nil {
			return err
		}
	case *machine.ConfigValidationErrorEvent:
		fields = append(fields, zap.Error(errors.New(event.GetError())))

		if err = a.handleConfigValidationFailedEvent(ctx, ip.String(), event); err != nil {
			return err
		}
	case *machine.ConfigLoadErrorEvent:
		fields = append(fields, zap.Error(errors.New(event.GetError())))

		if err = a.handleConfigLoadFailedEvent(ctx, ip.String(), event); err != nil {
			return err
		}
	case *machine.PhaseEvent:
		fields = append(fields, zap.String("phase", event.GetPhase()), zap.String("action", event.GetAction().String()))

		if err = a.handlePhaseEvent(ctx, ip.String(), event); err != nil {
			return err
		}
	case *machine.TaskEvent:
		fields = append(fields, zap.String("task", event.GetTask()), zap.String("action", event.GetAction().String()))
	case *machine.ServiceStateEvent:
		message = "service " + event.GetMessage()
		fields = append(fields, zap.String("service", event.GetService()), zap.String("action", event.GetAction().String()))
	case *machine.SequenceEvent:
		fields = append(fields, zap.String("sequence", event.GetSequence()), zap.String("action", event.GetAction().String()))

		if err = a.handleSequenceEvent(ctx, ip.String(), event); err != nil {
			return err
		}

		if event.GetError() != nil {
			err = errors.New(event.GetError().GetMessage())
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
	var err error

	switch {
	case event.GetSequence() == "install" && event.GetAction() == machine.SequenceEvent_NOOP && event.GetError().GetCode() == common.Code_FATAL:
		// this is ugly, but we need to handle this case
		if strings.Contains(event.GetError().GetMessage(), "unix.Reboot") {
			err = a.patchServerBinding(ctx, ip, func(serverbinding *sidero.ServerBinding) {
				conditions.MarkTrue(serverbinding, sidero.TalosInstalledCondition)
			})
		} else {
			err = a.patchServerBinding(ctx, ip, func(serverbinding *sidero.ServerBinding) {
				conditions.MarkFalse(serverbinding, sidero.TalosInstalledCondition, sidero.TalosInstallationFailedReason, clusterv1.ConditionSeverityError, "%s", event.GetError().GetMessage())
			})
		}
	case event.GetSequence() == "boot" && event.GetAction() == machine.SequenceEvent_START:
		err = a.patchServerBinding(ctx, ip, func(serverbinding *sidero.ServerBinding) {
			// if Talos reached 'boot' sequence, everything is good
			conditions.MarkTrue(serverbinding, sidero.TalosInstalledCondition)
			conditions.MarkTrue(serverbinding, sidero.TalosConfigValidatedCondition)
			conditions.MarkTrue(serverbinding, sidero.TalosConfigLoadedCondition)
		})
	}

	return err
}

func (a *Adapter) handleConfigLoadFailedEvent(ctx context.Context, ip string, event *machine.ConfigLoadErrorEvent) error {
	return a.patchServerBinding(ctx, ip, func(serverbinding *sidero.ServerBinding) {
		conditions.MarkFalse(serverbinding, sidero.TalosConfigLoadedCondition, sidero.TalosConfigLoadFailedReason, clusterv1.ConditionSeverityError, "%s", event.GetError())
	})
}

func (a *Adapter) handleConfigValidationFailedEvent(ctx context.Context, ip string, event *machine.ConfigValidationErrorEvent) error {
	return a.patchServerBinding(ctx, ip, func(serverbinding *sidero.ServerBinding) {
		conditions.MarkFalse(serverbinding, sidero.TalosConfigValidatedCondition, sidero.TalosConfigValidationFailedReason, clusterv1.ConditionSeverityError, "%s", event.GetError())
	})
}

func (a *Adapter) handlePhaseEvent(ctx context.Context, ip string, event *machine.PhaseEvent) error {
	var err error

	// handle only START events, STOP events come for both success and failures
	switch {
	case event.GetPhase() == "install" && event.GetAction() == machine.PhaseEvent_START:
		// starting phase install, mark as in progress
		err = a.patchServerBinding(ctx, ip, func(serverbinding *sidero.ServerBinding) {
			conditions.MarkFalse(serverbinding, sidero.TalosInstalledCondition, sidero.TalosInstallationInProgressReason, clusterv1.ConditionSeverityInfo, "")
		})
	case event.GetPhase() == "saveConfig" && event.GetAction() == machine.PhaseEvent_START:
		// if we reached this phase, config was validated and loaded
		err = a.patchServerBinding(ctx, ip, func(serverbinding *sidero.ServerBinding) {
			conditions.MarkTrue(serverbinding, sidero.TalosConfigValidatedCondition)
			conditions.MarkTrue(serverbinding, sidero.TalosConfigLoadedCondition)
		})
	}

	return err
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
