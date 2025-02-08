package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Redis struct {
	DSN string `json:"dsn"`
}

type Config struct {
	Port        int    `yaml:"port"`
	MetricsPort int    `yaml:"metricsPort"`
	Nats        string `yaml:"nats"`
	Redis       Redis  `yaml:"redis"`
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
