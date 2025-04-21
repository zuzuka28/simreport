package main

import (
	"github.com/zuzuka28/simreport/prj/shingleindex/internal/config"
	"github.com/zuzuka28/simreport/prj/shingleindex/pkg/cmd"

	"github.com/urfave/cli/v2"
)

func runAsyncProcessingCommand(c *cli.Context) error {
	ctx := c.Context

	cfg := ctx.Value(contextKeyConfig).(*config.Config) //nolint:forcetypeassert

	return cmd.RunAsyncProcessing(ctx, cfg) //nolint:wrapcheck
}
