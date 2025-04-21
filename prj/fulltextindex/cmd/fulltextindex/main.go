package main

import (
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{ //nolint:exhaustruct
		Name:  "github.com/zuzuka28/simreport/prj/fulltextindex-service API",
		Usage: "github.com/zuzuka28/simreport/prj/fulltextindex-service API service",
		Flags: []cli.Flag{
			&cli.StringFlag{ //nolint:exhaustruct
				Name:  "config",
				Value: "./config.yml",
				Usage: "path to the config file",
			},
		},
		Before: doBefore,
		Commands: []*cli.Command{
			{
				Name:   "run-intapi",
				Usage:  "run fulltextindex nats api",
				Action: runNATSAPICommand,
			},
			{
				Name:   "run-async-processing",
				Usage:  "run fulltextindex parse api",
				Action: runAsyncProcessingCommand,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("can't run application: " + err.Error())
	}
}
