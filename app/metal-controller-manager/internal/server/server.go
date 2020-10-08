// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"

	metalv1alpha1 "github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/app/metal-controller-manager/internal/api"
	"github.com/talos-systems/sidero/app/metal-controller-manager/pkg/client"
)

const (
	Port = "50100"
)

type server struct {
	api.UnimplementedAgentServer
	autoAccept bool
}

// CreateServer implements api.AgentServer.
func (s *server) CreateServer(ctx context.Context, in *api.CreateServerRequest) (*api.CreateServerResponse, error) {
	var config *rest.Config

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	c, err := client.NewClient(config)
	if err != nil {
		return nil, err
	}

	obj := &metalv1alpha1.Server{}

	if err = c.Get(ctx, types.NamespacedName{Name: in.GetSystemInformation().GetUuid()}, obj); err != nil {
		if !apierrors.IsNotFound(err) {
			return nil, err
		}

		obj = &metalv1alpha1.Server{
			TypeMeta: v1.TypeMeta{
				Kind:       "Server",
				APIVersion: metalv1alpha1.GroupVersion.Version,
			},
			ObjectMeta: v1.ObjectMeta{
				Name: in.GetSystemInformation().GetUuid(),
			},
			Spec: metalv1alpha1.ServerSpec{
				SystemInformation: &metalv1alpha1.SystemInformation{
					Manufacturer: in.GetSystemInformation().GetManufacturer(),
					ProductName:  in.GetSystemInformation().GetProductName(),
					Version:      in.GetSystemInformation().GetVersion(),
					SerialNumber: in.GetSystemInformation().GetSerialNumber(),
					SKUNumber:    in.GetSystemInformation().GetSkuNumber(),
					Family:       in.GetSystemInformation().GetFamily(),
				},
				CPU: &metalv1alpha1.CPUInformation{
					Manufacturer: in.GetCpu().GetManufacturer(),
					Version:      in.GetCpu().GetVersion(),
				},
				Accepted: s.autoAccept,
			},
		}

		if err = c.Create(ctx, obj); err != nil {
			return nil, err
		}

		log.Printf("Added %s", in.GetSystemInformation().GetUuid())
	}

	resp := &api.CreateServerResponse{}

	// Only return a wipe directive is the server is not clean *AND* it has been accepted.
	// This avoids the possibility of a random device PXE booting against us, registering, then getting blown away.
	if !obj.Status.IsClean && obj.Spec.Accepted {
		log.Printf("Server %q needs wipe", obj.Name)

		resp.Wipe = true
	}

	return resp, nil
}

// MarkServerAsWiped implements api.AgentServer.
func (s *server) MarkServerAsWiped(ctx context.Context, in *api.MarkServerAsWipedRequest) (*api.MarkServerAsWipedResponse, error) {
	var config *rest.Config

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	c, err := client.NewClient(config)
	if err != nil {
		return nil, err
	}

	obj := &metalv1alpha1.Server{}

	if err = c.Get(ctx, types.NamespacedName{Name: in.GetUuid()}, obj); err != nil {
		return nil, err
	}

	obj.Status.IsClean = true

	if err := c.Status().Update(ctx, obj); err != nil {
		return nil, err
	}

	resp := &api.MarkServerAsWipedResponse{}

	return resp, nil
}

func Serve(autoAccept bool) error {
	lis, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	api.RegisterAgentServer(s, &server{autoAccept: autoAccept})

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
