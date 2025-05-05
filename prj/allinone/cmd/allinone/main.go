package main

import (
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{ //nolint:exhaustruct
		Name:  "github.com/zuzuka28/simreport/prj/allinone",
		Usage: "github.com/zuzuka28/simreport/prj/allinone",
		Flags: []cli.Flag{
			&cli.StringFlag{ //nolint:exhaustruct
				Name:  "config",
				Value: "./config.yml",
				Usage: "path to the config file",
			},
		},
		Action: runApp,
		Before: doBefore,
		Commands: []*cli.Command{
			{
				Name:     "run-document-api",
				Usage:    "run document rest api",
				Category: "document",
				Action:   runDocumentRESTAPICommand,
			},
			{
				Name:     "run-document-intapi",
				Usage:    "run document nats api",
				Category: "document",
				Action:   runDocumentNATSAPICommand,
			},
			{
				Name:     "run-document-processing",
				Usage:    "run document processing",
				Category: "document",
				Action:   runDocumentAsyncProcessingCommand,
			},
			{
				Name:     "run-similarity-api",
				Usage:    "run similarity rest api",
				Category: "similarity",
				Action:   runSimilarityRESTAPICommand,
			},
			{
				Name:     "run-similarity-intapi",
				Usage:    "run similarity nats api",
				Category: "similarity",
				Action:   runSimilarityNATSAPICommand,
			},
			{
				Name:     "run-shingleindex-intapi",
				Usage:    "run shingleindex nats api",
				Category: "shingleindex",
				Action:   runShingleIndexNATSAPICommand,
			},
			{
				Name:     "run-shingleindex-processing",
				Usage:    "run shingleindex processing",
				Category: "shingleindex",
				Action:   runShingleIndexAsyncProcessingCommand,
			},
			{
				Name:     "run-fulltextindex-intapi",
				Usage:    "run fulltextindex nats api",
				Category: "fulltextindex",
				Action:   runFulltextIndexNATSAPICommand,
			},
			{
				Name:     "run-fulltextindex-processing",
				Usage:    "run fulltextindex processing",
				Category: "fulltextindex",
				Action:   runFulltextIndexAsyncProcessingCommand,
			},
			{
				Name:     "run-semanticindex-intapi",
				Usage:    "run semanticindex nats api",
				Category: "semanticindex",
				Action:   runSemanticIndexNATSAPICommand,
			},
			{
				Name:     "run-semanticindex-processing",
				Usage:    "run semanticindex processing",
				Category: "semanticindex",
				Action:   runSemanticIndexAsyncProcessingCommand,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("can't run application: " + err.Error())
	}
}
