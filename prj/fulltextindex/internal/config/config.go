package config

import (
	"fmt"

	"github.com/zuzuka28/simreport/lib/elasticutil"
	"github.com/zuzuka28/simreport/prj/fulltextindex/internal/repository/fulltextindex"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Nats         string             `yaml:"nats"`
	Elastic      elasticutil.Config `yaml:"elastic"`
	FulltextRepo fulltextindex.Opts `yaml:"fulltextRepo"`
}

func New(path string) (*Config, error) {
	cfg := new(Config)

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	return cfg, nil
}
