package config

import (
	"fmt"

	"github.com/zuzuka28/simreport/lib/elasticutil"
	"github.com/zuzuka28/simreport/lib/minioutil"
	"github.com/zuzuka28/simreport/prj/semanticindex/internal/repository/semanticindex"
	"github.com/zuzuka28/simreport/prj/semanticindex/internal/repository/vectorizer"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	MetricsPort    int                `yaml:"metricsPort"`
	Nats           string             `yaml:"nats"`
	Elastic        elasticutil.Config `yaml:"elastic"`
	S3             minioutil.Config   `yaml:"s3"`
	SemanticRepo   semanticindex.Opts `yaml:"semanticRepo"`
	VectorizerRepo vectorizer.Opts    `yaml:"vectorizerRepo"`
}

func New(path string) (*Config, error) {
	cfg := new(Config)

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	if cfg.MetricsPort == 0 {
		cfg.MetricsPort = 9000
	}

	return cfg, nil
}
