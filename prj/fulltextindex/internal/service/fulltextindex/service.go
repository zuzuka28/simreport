package fulltextindex

import (
	"context"
	"fmt"
	"fulltextindex/internal/model"
)

type Service struct {
	r Repository
}

func NewService(
	r Repository,
) *Service {
	return &Service{
		r: r,
	}
}

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

func (s *Service) Save(
	ctx context.Context,
	cmd model.DocumentSaveCommand,
) error {
	if err := s.r.Save(ctx, cmd); err != nil {
		return fmt.Errorf("save analyzed document: %w", err)
	}

	return nil
}
