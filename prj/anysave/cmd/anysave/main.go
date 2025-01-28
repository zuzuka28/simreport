package main

import (
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{ //nolint:exhaustruct
		Name:  "github.com/zuzuka28/simreport/prj/anysave-service API",
		Usage: "github.com/zuzuka28/simreport/prj/anysave-service API service",
		Flags: []cli.Flag{
			&cli.StringFlag{ //nolint:exhaustruct
				Name:  "config",
				Value: "./config.yml",
				Usage: "path to the config file",
			},
		},
		Action: runRestServer,
		Commands: []*cli.Command{
			{
				Name:   "run-api",
				Usage:  "run anysave rest api",
				Action: runRestServer,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("can't run application: " + err.Error())
	}
}
