package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/zuzuka28/simreport/lib/elasticutil"
	"github.com/zuzuka28/simreport/prj/tgbot/internal/bot"
	"github.com/zuzuka28/simreport/prj/tgbot/internal/repository/userstate"
)

type Config struct {
	Bot           bot.Config         `yaml:"bot"`
	Nats          string             `yaml:"nats"`
	Elastic       elasticutil.Config `yaml:"elastic"`
	UserStateRepo userstate.Config   `yaml:"userStateRepo"`
}

func New(path string) (*Config, error) {
	cfg := new(Config)

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	return cfg, nil
}
