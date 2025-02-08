package config

import (
	"fmt"

	"github.com/zuzuka28/simreport/prj/similarity/internal/repository/analyzehistory"

	"github.com/zuzuka28/simreport/lib/elasticutil"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port               int                 `yaml:"port"`
	MetricsPort        int                 `yaml:"metricsPort"`
	Nats               string              `yaml:"nats"`
	Elastic            elasticutil.Config  `yaml:"elastic"`
	AnalyzeHistoryRepo analyzehistory.Opts `yaml:"analyzeHistoryRepo"`
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
