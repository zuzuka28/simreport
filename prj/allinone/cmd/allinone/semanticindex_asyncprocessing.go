package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/zuzuka28/simreport/prj/allinone/internal/config"
	"github.com/zuzuka28/simreport/prj/semanticindex/pkg/cmd"

	"github.com/urfave/cli/v2"
)

func runSemanticIndexAsyncProcessingCommand(c *cli.Context) error {
	ctx := c.Context

	cfg := ctx.Value(contextKeyConfig).(*config.Config)         //nolint:forcetypeassert
	reg := ctx.Value(contextKeyRegistry).(*prometheus.Registry) //nolint:forcetypeassert

	return cmd.RunAsyncProcessing(
		ctx,
		&cfg.SemanticIndexService,
		cmd.WithPrometheusRegistrer(reg),
	) //nolint:wrapcheck
}
