package semanticindex

import (
	"github.com/elastic/go-elasticsearch/v8"
)

type Opts struct {
	Index string `yaml:"index"`
}

type Repository struct {
	cli   *elasticsearch.Client
	index string
}

func NewRepository(
	opts Opts,
	es *elasticsearch.Client,
) (*Repository, error) {
	return &Repository{
		cli:   es,
		index: opts.Index,
	}, nil
}
