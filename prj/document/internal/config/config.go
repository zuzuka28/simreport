package config

import (
	"fmt"

	"github.com/zuzuka28/simreport/prj/document/internal/repository/attribute"
	"github.com/zuzuka28/simreport/prj/document/internal/repository/document"

	"github.com/zuzuka28/simreport/lib/elasticutil"
	"github.com/zuzuka28/simreport/lib/minioutil"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port          int                `yaml:"port"`
	S3            minioutil.Config   `yaml:"s3"`
	Nats          string             `yaml:"nats"`
	Tika          string             `yaml:"tika"`
	Elastic       elasticutil.Config `yaml:"elastic"`
	DocumentRepo  document.Opts      `yaml:"documentRepo"`
	AttributeRepo attribute.Opts     `yaml:"attributeRepo"`
}

func New(path string) (*Config, error) {
	cfg := new(Config)

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	return cfg, nil
}
