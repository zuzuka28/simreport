package shingleindex

import (
	"hash/fnv"

	"github.com/redis/go-redis/v9"
	"github.com/zuzuka28/simreport/lib/minhashlsh"
	"github.com/zuzuka28/simreport/lib/minhashlsh/redisstorage"
)

const (
	prefix       = "shingleindex"
	permutations = 512
	bands        = 64
)

var hasher = func(b []byte) uint64 {
	h := fnv.New64a()
	h.Write(b)
	return h.Sum64()
}

type Opts struct{}

type Repository struct {
	lsh *minhashlsh.MinhashLSH
}

func NewRepository(
	opts Opts,
	cli *redis.Client,
) (*Repository, error) {
	lsh := minhashlsh.New(
		redisstorage.New(cli),
		prefix,
		permutations,
		bands,
		hasher,
	)

	return &Repository{
		lsh: lsh,
	}, nil
}
