package shingleindex

import (
	"context"
	"fmt"
	"hash/fnv"
	"shingleindex/internal/model"

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

func (r *Repository) Save(
	ctx context.Context,
	cmd model.MinhashSaveCommand,
) error {
	if err := r.lsh.Insert(ctx, cmd.DocumentID, cmd.Minhash); err != nil {
		return fmt.Errorf("insert into lsh: %w", err)
	}

	return nil
}

func (r *Repository) SearchSimilar(
	ctx context.Context,
	query model.MinhashSimilarQuery,
) ([]*model.MinhashSimilarMatch, error) {
	candidates, err := r.lsh.Query(ctx, query.Minhash)
	if err != nil {
		return nil, fmt.Errorf("query candidates lsh: %w", err)
	}

	return mapCandidatesToMatches(candidates), nil
}
