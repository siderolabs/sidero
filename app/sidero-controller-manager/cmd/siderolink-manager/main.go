// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/netip"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/siderolabs/siderolink/api/siderolink"
	"github.com/siderolabs/siderolink/pkg/wireguard"

	"github.com/siderolabs/sidero/app/sidero-controller-manager/internal/siderolink"
	"github.com/siderolabs/sidero/app/sidero-controller-manager/pkg/constants"
)

var (
	wireguardEndpoint string
	wireguardPort     int
)

func main() {
	flag.StringVar(&wireguardEndpoint, "wireguard-endpoint", "", "The endpoint (IP address) SideroLink can be reached at from the servers.")
	flag.IntVar(&wireguardPort, "wireguard-port", 51821, "The TCP port SideroLink can be reached at from the servers.")

	flag.Parse()

	if wireguardEndpoint == "-" {
		wireguardEndpoint = ""
	}

	if wireguardEndpoint == "" {
		if endpoint, ok := os.LookupEnv("API_ENDPOINT"); ok {
			wireguardEndpoint = endpoint
		} else {
			log.Fatal("no Wireguard endpoint found")
		}
	}

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}
}

func recoveryHandler(logger *zap.Logger) grpc_recovery.RecoveryHandlerFunc {
	return func(p interface{}) error {
		if logger != nil {
			logger.Error("grpc panic", zap.Any("panic", p), zap.Stack("stack"))
		}

		return status.Errorf(codes.Internal, "%v", p)
	}
}

func run() error {
	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	zap.ReplaceGlobals(logger)
	zap.RedirectStdLog(logger)

	metalclient, kubeconfig, err := getMetalClient()
	if err != nil {
		return fmt.Errorf("error building runtime client: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	lis, err := net.Listen("tcp", constants.SideroLinkInternalAPIEndpoint)
	if err != nil {
		return fmt.Errorf("error listening for gRPC API: %w", err)
	}

	wireguardEndpointIP, err := getIPForHost(wireguardEndpoint)
	if err != nil {
		return err
	}

	if wireguardEndpoint != wireguardEndpointIP {
		logger.Sugar().Infof("resolved wireguard endpoint %s to %s", wireguardEndpoint, wireguardEndpointIP)
	}

	siderolink.Cfg.WireguardEndpoint = net.JoinHostPort(wireguardEndpointIP, strconv.Itoa(wireguardPort))

	if err = siderolink.Cfg.LoadOrCreate(ctx, metalclient); err != nil {
		return err
	}

	wireguardEndpoint, err := netip.ParseAddrPort(siderolink.Cfg.WireguardEndpoint)
	if err != nil {
		return fmt.Errorf("invalid Wireguard endpoint: %w", err)
	}

	wgDevice, err := wireguard.NewDevice(siderolink.Cfg.ServerAddress, siderolink.Cfg.PrivateKey, wireguardEndpoint.Port(), false, logger)
	if err != nil {
		return fmt.Errorf("error initializing wgDevice: %w", err)
	}

	defer wgDevice.Close() //nolint:errcheck

	grpc_zap.ReplaceGrpcLoggerV2(logger)

	recoveryOpt := grpc_recovery.WithRecoveryHandler(recoveryHandler(logger))

	serverOptions := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_recovery.UnaryServerInterceptor(recoveryOpt),
		),
		grpc.ChainStreamInterceptor(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(logger),
			grpc_recovery.StreamServerInterceptor(recoveryOpt),
		),
	}

	srv := siderolink.NewServer(&siderolink.Cfg, metalclient)

	peers := siderolink.NewPeerState(kubeconfig, logger)

	s := grpc.NewServer(serverOptions...)
	pb.RegisterProvisionServiceServer(s, srv)

	eg.Go(func() error {
		return wgDevice.Run(ctx, logger, peers)
	})

	eg.Go(func() error {
		return peers.Run(ctx)
	})

	eg.Go(func() error {
		return s.Serve(lis)
	})

	eg.Go(func() error {
		<-ctx.Done()

		s.Stop()

		return nil
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, grpc.ErrServerStopped) && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

func getIPForHost(host string) (string, error) {
	parsedIP, err := netip.ParseAddr(host)
	if err == nil {
		return parsedIP.String(), nil
	}

	resolvedIPs, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}

	if len(resolvedIPs) == 0 {
		return "", fmt.Errorf("no IPs found for %s", host)
	}

	return resolvedIPs[0].String(), nil
}
