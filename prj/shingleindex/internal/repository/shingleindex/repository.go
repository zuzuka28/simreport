package shingleindex

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/zuzuka28/simreport/lib/minhashlsh"
	"github.com/zuzuka28/simreport/lib/minhashlsh/elasticstorage"
)

type Opts elasticstorage.Config

type Repository struct {
	lsh *minhashlsh.MinhashLSH
}

func NewRepository(
	opts Opts,
	cli *elasticsearch.Client,
) (*Repository, error) {
	lsh := minhashlsh.New(
		elasticstorage.New(elasticstorage.Config(opts), cli),
		prefix,
		permutations,
		bands,
		hasher,
	)

	return &Repository{
		lsh: lsh,
	}, nil
}
