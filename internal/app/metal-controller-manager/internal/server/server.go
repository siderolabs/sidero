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
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	metal1alpha1 "github.com/talos-systems/sidero/internal/app/metal-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/internal/app/metal-controller-manager/internal/api"
	"github.com/talos-systems/sidero/internal/app/metal-controller-manager/pkg/client"
)

const (
	Port = "50100"
)

type server struct {
	api.UnimplementedDiscoveryServer
}

// CreateServer implements api.DiscoveryServer.
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

	obj := &metal1alpha1.Server{
		TypeMeta: v1.TypeMeta{
			Kind:       "Server",
			APIVersion: metal1alpha1.GroupVersion.Version,
		},
		ObjectMeta: v1.ObjectMeta{
			Name: in.GetSystemInformation().GetUuid(),
		},
	}

	_, err = controllerutil.CreateOrUpdate(context.Background(), c, obj, func() error {
		obj.Spec = metal1alpha1.ServerSpec{
			SystemInformation: &metal1alpha1.SystemInformation{
				Manufacturer: in.GetSystemInformation().GetManufacturer(),
				ProductName:  in.GetSystemInformation().GetProductName(),
				Version:      in.GetSystemInformation().GetVersion(),
				SerialNumber: in.GetSystemInformation().GetSerialNumber(),
				SKUNumber:    in.GetSystemInformation().GetSkuNumber(),
				Family:       in.GetSystemInformation().GetFamily(),
			},
			CPU: &metal1alpha1.CPUInformation{
				Manufacturer: in.GetCpu().GetManufacturer(),
				Version:      in.GetCpu().GetVersion(),
			},
		}

		return nil
	})
	if err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return nil, err
		}

		log.Printf("%s already exists", in.GetSystemInformation().GetUuid())
	} else {
		log.Printf("Added %s", in.GetSystemInformation().GetUuid())
	}

	return &api.CreateServerResponse{}, nil
}

// Server starts the server.
func Serve() error {
	lis, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	api.RegisterDiscoveryServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
