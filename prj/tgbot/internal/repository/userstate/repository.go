package userstate

import (
	"github.com/elastic/go-elasticsearch/v8"
)

type Config struct {
	Index string `yaml:"index"`
}

type Repository struct {
	es    *elasticsearch.Client
	index string

	m Metrics
}

func NewRepository(
	cfg Config,
	cli *elasticsearch.Client,
	m Metrics,
) *Repository {
	return &Repository{
		es:    cli,
		index: cfg.Index,
		m:     m,
	}
}
