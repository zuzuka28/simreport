package document

import (
	"context"
	"document/internal/model"
	"fmt"
)

func (s *Service) Search(
	ctx context.Context,
	query model.DocumentSearchQuery,
) ([]model.Document, error) {
	res, err := s.r.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("search documents: %w", err)
	}

	return res, nil
}
