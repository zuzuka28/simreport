package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	cli "github.com/urfave/cli/v2"
	"github.com/zuzuka28/simreport/prj/semanticindex/internal/config"
	"github.com/zuzuka28/simreport/prj/semanticindex/internal/metrics"
)

type contextKey int

const (
	contextKeyConfig contextKey = iota + 1
	contextKeyRegistry
)

func doBefore(c *cli.Context) error {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := config.New(c.String("config"))
	if err != nil {
		return fmt.Errorf("new config: %w", err)
	}

	c.Context = context.WithValue(c.Context, contextKeyConfig, cfg)

	reg := prometheus.NewRegistry()

	c.Context = context.WithValue(c.Context, contextKeyRegistry, reg)

	gcol := collectors.NewGoCollector()
	pcol := collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}) //nolint:exhaustruct

	reg.MustRegister(gcol, pcol)

	msrv := metrics.NewMetricsServer(cfg.MetricsPort, reg)

	go func() {
		_ = msrv.Start(c.Context)
	}()

	return nil
}
