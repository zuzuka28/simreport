package shingleindex

import (
	"context"
	"fmt"
	"shingleindex/internal/model"
)

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
