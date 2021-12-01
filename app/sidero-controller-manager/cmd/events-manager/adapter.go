// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/talos-systems/siderolink/pkg/events"

	"github.com/talos-systems/talos/pkg/machinery/api/machine"
)

// Adapter implents gRPC API.
type Adapter struct {
	Sink *events.Sink

	logger      *zap.Logger
	metalClient runtimeclient.Client
}

// NewAdapter initializes new server.
func NewAdapter(metalClient runtimeclient.Client, logger *zap.Logger) *Adapter {
	return &Adapter{
		logger:      logger,
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

	switch event := event.Payload.(type) {
	case *machine.AddressEvent:
		fields = append(fields, zap.String("hostname", event.GetHostname()), zap.String("addresses", strings.Join(event.GetAddresses(), ",")))
	case *machine.ConfigValidationErrorEvent:
		fields = append(fields, zap.Error(fmt.Errorf(event.GetError())))
	case *machine.ConfigLoadErrorEvent:
		fields = append(fields, zap.Error(fmt.Errorf(event.GetError())))
	case *machine.PhaseEvent:
		fields = append(fields, zap.String("phase", event.GetPhase()), zap.String("action", event.GetAction().String()))
	case *machine.TaskEvent:
		fields = append(fields, zap.String("task", event.GetTask()), zap.String("action", event.GetAction().String()))
	case *machine.ServiceStateEvent:
		message = "service " + event.GetMessage()
		fields = append(fields, zap.String("service", event.GetService()), zap.String("action", event.GetAction().String()))
	case *machine.SequenceEvent:
		fields = append(fields, zap.String("sequence", event.GetSequence()), zap.String("action", event.GetAction().String()))

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
