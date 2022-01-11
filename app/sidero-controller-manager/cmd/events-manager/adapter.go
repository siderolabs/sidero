// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	sidero "github.com/talos-systems/sidero/app/caps-controller-manager/api/v1alpha3"

	"github.com/talos-systems/siderolink/pkg/events"

	"github.com/talos-systems/talos/pkg/machinery/api/machine"
)

// Adapter implents gRPC API.
type Adapter struct {
	Sink *events.Sink

	logger      *zap.Logger
	metalClient runtimeclient.Client
	kubeconfig  *rest.Config
	nodesMu     sync.Mutex
	nodes       map[string]types.NamespacedName
}

// NewAdapter initializes new server.
func NewAdapter(metalClient runtimeclient.Client, kubeconfig *rest.Config, logger *zap.Logger) *Adapter {
	return &Adapter{
		logger:      logger,
		kubeconfig:  kubeconfig,
		metalClient: metalClient,
		nodes:       map[string]types.NamespacedName{},
	}
}

func (a *Adapter) Run(ctx context.Context) error {
	dc, err := dynamic.NewForConfig(a.kubeconfig)
	if err != nil {
		return err
	}

	// Create a factory object that can generate informers for resource types
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dc, 10*time.Minute, "", nil)

	informerFactory := factory.ForResource(sidero.GroupVersion.WithResource("serverbindings"))
	informer := informerFactory.Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(new interface{}) { a.notify(nil, new) },
		UpdateFunc: a.notify,
		DeleteFunc: func(old interface{}) { a.notify(old, nil) },
	})

	informer.Run(ctx.Done())

	return nil
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

	parts := strings.Split(event.Node, ":")
	ip := strings.Join(parts[:len(parts)-1], ":")

	switch event := event.Payload.(type) {
	case *machine.AddressEvent:
		fields = append(fields, zap.String("hostname", event.GetHostname()), zap.String("addresses", strings.Join(event.GetAddresses(), ",")))

		err = a.patchServerBinding(ctx, ip, func(serverbinding *sidero.ServerBinding) {
			serverbinding.Spec.Addresses = event.Addresses
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
	if event.GetSequence() == "install" &&
		event.GetAction() == machine.SequenceEvent_STOP {
		var callback func(*sidero.ServerBinding)

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

		if e := a.patchServerBinding(ctx, ip, callback); e != nil {
			return e
		}
	}

	return nil
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
	a.nodesMu.Lock()
	defer a.nodesMu.Unlock()

	name, ok := a.nodes[ip]
	if !ok {
		return fmt.Errorf("failed to find ServerBindings for ip %s", ip)
	}

	var serverbinding sidero.ServerBinding
	if err := a.metalClient.Get(ctx, name, &serverbinding); err != nil {
		return err
	}

	patchHelper, err := patch.NewHelper(&serverbinding, a.metalClient)
	if err != nil {
		return err
	}

	callback(&serverbinding)

	return patchHelper.Patch(ctx, &serverbinding)
}

func (a *Adapter) notify(old, new interface{}) {
	var oldServerBinding, newServerBinding *sidero.ServerBinding

	if old != nil {
		oldServerBinding = &sidero.ServerBinding{}

		err := runtime.DefaultUnstructuredConverter.
			FromUnstructured(old.(*unstructured.Unstructured).UnstructuredContent(), oldServerBinding)
		if err != nil {
			a.logger.Error("failed converting old event object", zap.Error(err))

			return
		}
	}

	if new != nil {
		newServerBinding = &sidero.ServerBinding{}

		err := runtime.DefaultUnstructuredConverter.
			FromUnstructured(new.(*unstructured.Unstructured).UnstructuredContent(), newServerBinding)
		if err != nil {
			a.logger.Error("failed converting new event object", zap.Error(err))

			return
		}
	}

	a.nodesMu.Lock()
	defer a.nodesMu.Unlock()

	if new == nil {
		delete(a.nodes, oldServerBinding.Spec.SideroLink.NodeAddress)
	} else {
		address := newServerBinding.Spec.SideroLink.NodeAddress
		if address == "" {
			return
		}

		address = fmt.Sprintf("[%s]", strings.Split(address, "/")[0])

		if old != nil {
			delete(a.nodes, oldServerBinding.Spec.SideroLink.NodeAddress)
		}

		a.nodes[address] = types.NamespacedName{
			Name:      newServerBinding.GetName(),
			Namespace: newServerBinding.GetNamespace(),
		}

		a.logger.Info("new node mapping", zap.String("ip", newServerBinding.Spec.SideroLink.NodeAddress), zap.String("server", newServerBinding.GetName()))
	}
}
