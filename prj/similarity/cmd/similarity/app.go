package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/zuzuka28/simreport/prj/similarity/internal/provider"

	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

func runApp(c *cli.Context) error {
	cfg, err := provider.InitConfig(c.String("config"))
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	restapi, err := provider.InitRestAPI(c.Context, cfg)
	if err != nil {
		return fmt.Errorf("init rest api: %w", err)
	}

	natsapi, err := provider.InitNatsAPI(c.Context, cfg)
	if err != nil {
		return fmt.Errorf("init nats api: %w", err)
	}

	errCh := make(chan error)

	go func() {
		eg, egCtx := errgroup.WithContext(c.Context)

		eg.Go(func() error {
			return restapi.Start(egCtx)
		})

		eg.Go(func() error {
			return natsapi.Start(egCtx)
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
