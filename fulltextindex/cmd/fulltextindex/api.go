package main

import (
	"fmt"
	"fulltextindex/internal/provider"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
)

func runMicroServer(c *cli.Context) error {
	cfg, err := provider.InitConfig(c.String("config"))
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	api, err := provider.InitNatsMicroAPI(c.Context, cfg)
	if err != nil {
		return fmt.Errorf("init api: %w", err)
	}

	errCh := make(chan error)

	go func() {
		if err := api.Start(c.Context); err != nil {
			errCh <- fmt.Errorf("run webserver: %w", err)
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
