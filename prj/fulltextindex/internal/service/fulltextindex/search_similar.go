package fulltextindex

import (
	"context"
	"fmt"
	"fulltextindex/internal/model"
)

func (s *Service) SearchSimilar(
	ctx context.Context,
	query model.DocumentSimilarQuery,
) ([]model.DocumentSimilarMatch, error) {
	res, err := s.r.SearchSimilar(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("search fulltext similar: %w", err)
	}

	return res, nil
}
