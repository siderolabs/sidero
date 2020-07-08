/*
Copyright 2020 Talos Systems, Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"context"
	"fmt"
	"log"
	"net"

	metal1alpha1 "github.com/talos-systems/sidero/internal/app/metal-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/internal/app/metal-controller-manager/internal/api"
	"github.com/talos-systems/sidero/internal/app/metal-controller-manager/pkg/client"
	"google.golang.org/grpc"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	Port = "50100"
)

type server struct {
	api.UnimplementedDiscoveryServer
}

// CreateServer implements api.DiscoveryServer
func (s *server) CreateServer(ctx context.Context, in *api.CreateServerRequest) (*api.CreateServerResponse, error) {
	var (
		config *rest.Config
	)

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

func ServerAPI() error {
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
