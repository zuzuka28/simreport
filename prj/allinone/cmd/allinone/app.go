package main

import (
	"github.com/prometheus/client_golang/prometheus"
	cli "github.com/urfave/cli/v2"
	"github.com/zuzuka28/simreport/prj/allinone/internal/config"
	"github.com/zuzuka28/simreport/prj/allinone/pkg/cmd"
)

func runApp(c *cli.Context) error {
	ctx := c.Context

	cfg := ctx.Value(contextKeyConfig).(*config.Config)         //nolint:forcetypeassert
	reg := ctx.Value(contextKeyRegistry).(*prometheus.Registry) //nolint:forcetypeassert

	return cmd.RunApp(c.Context, cfg, cmd.WithPrometheusRegistrer(reg))
}
