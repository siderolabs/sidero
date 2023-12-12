// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package siderolink

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"net/netip"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/cluster-api/util/patch"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	pb "github.com/siderolabs/siderolink/api/siderolink"

	sidero "github.com/siderolabs/sidero/app/caps-controller-manager/api/v1alpha3"
)

// Server implements gRPC API.
type Server struct {
	pb.UnimplementedProvisionServiceServer

	cfg         *Config
	metalClient runtimeclient.Client
}

// NewServer initializes new server.
func NewServer(cfg *Config, metalClient runtimeclient.Client) *Server {
	return &Server{
		cfg:         cfg,
		metalClient: metalClient,
	}
}

// Provision the SideroLink for the server by UUID.
func (srv *Server) Provision(ctx context.Context, req *pb.ProvisionRequest) (*pb.ProvisionResponse, error) {
	var serverbinding sidero.ServerBinding

	if err := srv.metalClient.Get(ctx, types.NamespacedName{Name: req.NodeUuid}, &serverbinding); err != nil {
		if apierrors.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("server binding %q not found", req.NodeUuid))
		}

		return nil, err
	}

	patchHelper, err := patch.NewHelper(&serverbinding, srv.metalClient)
	if err != nil {
		return nil, err
	}

	var nodeAddress netip.Prefix

	if serverbinding.Spec.SideroLink.NodeAddress != "" {
		// find already provisioned address
		nodeAddress, err = netip.ParsePrefix(serverbinding.Spec.SideroLink.NodeAddress)
		if err != nil {
			return nil, err
		}
	} else {
		// generate random address for the node
		raw := srv.cfg.Subnet.Addr().As16()
		salt := make([]byte, 8)

		_, err := io.ReadFull(rand.Reader, salt)
		if err != nil {
			return nil, err
		}

		copy(raw[8:], salt)

		nodeAddress = netip.PrefixFrom(netip.AddrFrom16(raw), srv.cfg.Subnet.Bits())

		serverbinding.Spec.SideroLink.NodeAddress = nodeAddress.String()
	}

	pubKey, err := wgtypes.ParseKey(req.NodePublicKey)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("error parsing Wireguard key: %s", err))
	}

	serverbinding.Spec.SideroLink.NodePublicKey = pubKey.String()

	if err = patchHelper.Patch(ctx, &serverbinding); err != nil {
		return nil, err
	}

	return &pb.ProvisionResponse{
		ServerEndpoint:    []string{srv.cfg.WireguardEndpoint},
		ServerPublicKey:   srv.cfg.PublicKey.String(),
		ServerAddress:     srv.cfg.ServerAddress.Addr().String(),
		NodeAddressPrefix: nodeAddress.String(),
	}, nil
}
