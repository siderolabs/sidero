// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/siderolink"
	"github.com/talos-systems/siderolink/api/events"
	sink "github.com/talos-systems/siderolink/pkg/events"
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}
}

func run() error {
	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	zap.ReplaceGlobals(logger)
	zap.RedirectStdLog(logger)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	address := fmt.Sprintf(":%d", siderolink.EventsSinkPort)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("error listening for gRPC API: %w", err)
	}

	s := grpc.NewServer()

	client, kubeconfig, err := getMetalClient()
	if err != nil {
		return fmt.Errorf("error getting metal client: %w", err)
	}

	adapter := NewAdapter(client,
		kubeconfig,
		logger.With(zap.String("component", "sink")),
	)

	srv := sink.NewSink(adapter)

	events.RegisterEventSinkServiceServer(s, srv)

	eg.Go(func() error {
		return adapter.Run(ctx)
	})

	eg.Go(func() error {
		logger.Info("started gRPC event sink", zap.String("address", address))

		return s.Serve(lis)
	})

	eg.Go(func() error {
		<-ctx.Done()

		s.Stop()

		return nil
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, grpc.ErrServerStopped) && errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}
