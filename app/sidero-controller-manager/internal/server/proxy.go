// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package server

import (
	"context"
	"strings"
	"sync"

	"github.com/talos-systems/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/talos-systems/sidero/app/sidero-controller-manager/pkg/constants"
)

// director proxy passes gRPC APIs to sub-components based on API method name.
func director(ctx context.Context, fullMethodName string) (proxy.Mode, []proxy.Backend, error) {
	switch {
	case strings.HasPrefix(fullMethodName, "/sidero.link."):
		return proxy.One2One, []proxy.Backend{sideroLinkAPI}, nil
	default:
		return proxy.One2One, nil, status.Errorf(codes.Unimplemented, "Unknown method")
	}
}

// backend performs proxying to another Sidero component.
type backend struct {
	target string

	mu   sync.Mutex
	conn *grpc.ClientConn
}

func (b *backend) String() string {
	return b.target
}

// GetConnection returns a grpc connection to the backend.
func (b *backend) GetConnection(ctx context.Context) (context.Context, *grpc.ClientConn, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	outCtx := metadata.NewOutgoingContext(ctx, md)

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.conn != nil {
		return outCtx, b.conn, nil
	}

	var err error
	b.conn, err = grpc.DialContext(
		ctx,
		b.target,
		grpc.WithInsecure(),
		grpc.WithCodec(proxy.Codec()), //nolint:staticcheck
	)

	return outCtx, b.conn, err
}

// AppendInfo is called to enhance response from the backend with additional data.
func (b *backend) AppendInfo(streaming bool, resp []byte) ([]byte, error) {
	return resp, nil
}

// BuildError is called to convert error from upstream into response field.
func (b *backend) BuildError(streaming bool, err error) ([]byte, error) {
	return nil, err
}

var sideroLinkAPI = &backend{target: constants.SideroLinkInternalAPIEndpoint}
