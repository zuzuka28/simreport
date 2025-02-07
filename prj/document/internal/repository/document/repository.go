package document

import (
	"github.com/elastic/go-elasticsearch/v8"
)

type Opts struct {
	Index string `yaml:"index"`
}

type Repository struct {
	cli   *elasticsearch.Client
	index string

	m Metrics
}

func NewRepository(
	opts Opts,
	es *elasticsearch.Client,
	m Metrics,
) (*Repository, error) {
	return &Repository{
		cli:   es,
		index: opts.Index,
		m:     m,
	}, nil
}
