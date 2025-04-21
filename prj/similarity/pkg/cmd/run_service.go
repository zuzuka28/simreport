package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type runnable interface {
	Start(ctx context.Context) error
}

func runService(
	ctx context.Context,
	svc runnable,
	shutdown func(),
) error {
	errCh := make(chan error)

	doneCh := make(chan struct{})

	go func() {
		if err := svc.Start(ctx); err != nil {
			errCh <- fmt.Errorf("run service: %w", err)
		}

		doneCh <- struct{}{}
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-doneCh:
		slog.Info("graceful shutdown initiated", "status", "done")

		shutdown()

		return nil

	case err := <-errCh:
		slog.Info("graceful shutdown initiated", "err", err.Error())

		shutdown()

		return err

	case sig := <-osSignals:
		slog.Info("graceful shutdown initiated", "sig", sig)

		shutdown()

		return nil
	}
}
