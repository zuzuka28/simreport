package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/zuzuka28/simreport/prj/similarity/internal/config"
	"github.com/zuzuka28/simreport/prj/similarity/pkg/cmd"

	"github.com/urfave/cli/v2"
)

func runRESTAPICommand(c *cli.Context) error {
	ctx := c.Context

	cfg := ctx.Value(contextKeyConfig).(*config.Config)         //nolint:forcetypeassert
	reg := ctx.Value(contextKeyRegistry).(*prometheus.Registry) //nolint:forcetypeassert

	return cmd.RunRESTAPI(ctx, cfg, cmd.WithPrometheusRegistrer(reg)) //nolint:wrapcheck
}
