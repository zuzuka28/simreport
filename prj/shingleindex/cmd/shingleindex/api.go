package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/zuzuka28/simreport/prj/shingleindex/internal/provider"

	"github.com/urfave/cli/v2"
)

func runMicroServer(c *cli.Context) error {
	cfg, err := provider.InitConfig(c.String("config"))
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	_, err = provider.InitNatsMicroAPI(c.Context, cfg)
	if err != nil {
		return fmt.Errorf("init api: %w", err)
	}

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	sig := <-osSignals
	slog.Warn("got signal", "sig", sig)

	// TODO: Add graceful shutdown logic here

	return nil
}
