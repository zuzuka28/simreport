package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"simrep/internal/provider"
	"syscall"

	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

func runApp(c *cli.Context) error {
	cfg, err := provider.InitConfig(c.String("config"))
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	api, err := provider.InitRestAPI(c.Context, cfg)
	if err != nil {
		return fmt.Errorf("init api: %w", err)
	}

	processing, err := provider.InitDocumentPipeline(c.Context, cfg)
	if err != nil {
		return fmt.Errorf("init api: %w", err)
	}

	errCh := make(chan error)

	go func() {
		eg, egCtx := errgroup.WithContext(c.Context)

		eg.Go(func() error {
			return api.Start(egCtx)
		})

		eg.Go(func() error {
			return processing.Start(c.Context)
		})

		if err := eg.Wait(); err != nil {
			errCh <- fmt.Errorf("run app: %w", err)
		}
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err

	case sig := <-osSignals:
		slog.Warn("got signal", "sig", sig)

		// TODO: Add graceful shutdown logic here

		return nil
	}
}
