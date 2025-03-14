package fulltextindex

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/similarity/internal/model"
)

func (s *Service) SearchSimilar(
	ctx context.Context,
	query model.SimilarityQuery,
) ([]*model.SimilarityMatch, error) {
	res, err := s.r.SearchSimilar(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("search fulltext similar: %w", err)
	}

	return res, nil
}
