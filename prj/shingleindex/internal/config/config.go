package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Redis struct {
	DSN string `json:"dsn"`
}

type Config struct {
	Port  int    `yaml:"port"`
	Nats  string `yaml:"nats"`
	Redis Redis  `yaml:"redis"`
}

func New(path string) (*Config, error) {
	cfg := new(Config)

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	return cfg, nil
}
