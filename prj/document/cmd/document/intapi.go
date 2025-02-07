package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/zuzuka28/simreport/prj/document/internal/metrics"
	"github.com/zuzuka28/simreport/prj/document/internal/provider"

	"github.com/urfave/cli/v2"
)

func runNatsServer(c *cli.Context) error {
	cfg, err := provider.InitConfig(c.String("config"))
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	api, err := provider.InitNatsAPI(c.Context, cfg)
	if err != nil {
		return fmt.Errorf("init api: %w", err)
	}

	reg := metrics.NewRegistry()
	reg.MustRegister(provider.ProvideMetrics().Collectors()...)
	metricsserv := metrics.NewMetricsServer(cfg.MetricsPort, reg)

	errCh := make(chan error)

	go func() {
		if err := metricsserv.Start(); err != nil {
			errCh <- fmt.Errorf("run metrics server: %w", err)
		}
	}()

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
