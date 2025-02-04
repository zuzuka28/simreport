package main

import (
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{ //nolint:exhaustruct
		Name:  "github.com/zuzuka28/simreport/prj/semanticindex-service API",
		Usage: "github.com/zuzuka28/simreport/prj/semanticindex-service API service",
		Flags: []cli.Flag{
			&cli.StringFlag{ //nolint:exhaustruct
				Name:  "config",
				Value: "./config.yml",
				Usage: "path to the config file",
			},
		},
		Action: runApp,
		Commands: []*cli.Command{
			{
				Name:   "run-api",
				Usage:  "run semanticindex nats api",
				Action: runMicroServer,
			},
			{
				Name:   "run-async-processing",
				Usage:  "run semantictindex parse api",
				Action: runAsyncProcessing,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("can't run application: " + err.Error())
	}
}
