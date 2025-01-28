package shingleindex

import (
	"context"
	"document/internal/model"
	"fmt"
)

func (s *Service) SearchSimilar(
	ctx context.Context,
	query model.DocumentSimilarQuery,
) ([]*model.DocumentSimilarMatch, error) {
	res, err := s.r.SearchSimilar(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("search shingle similar: %w", err)
	}

	return res, nil
}
