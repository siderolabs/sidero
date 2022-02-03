// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/siderolink"
	"github.com/talos-systems/siderolink/pkg/logreceiver"
)

func main() {
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

	metalclient, kubeconfig, err := getMetalClient()
	if err != nil {
		return fmt.Errorf("error building runtime client: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", siderolink.LogReceiverPort))
	if err != nil {
		return fmt.Errorf("error listening for log endpoint: %w", err)
	}

	annotator := siderolink.NewAnnotator(metalclient, kubeconfig, logger)

	srv, err := logreceiver.NewServer(logger, listener, logHandler(logger, annotator))
	if err != nil {
		return fmt.Errorf("error initializing log receiver: %w", err)
	}

	eg.Go(func() error {
		return annotator.Run(ctx)
	})

	eg.Go(func() error {
		return srv.Serve()
	})

	eg.Go(func() error {
		<-ctx.Done()

		srv.Stop()

		return nil
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}
