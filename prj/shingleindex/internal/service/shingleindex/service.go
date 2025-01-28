package shingleindex

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/shingleindex/internal/model"
)

type Service struct {
	search *searchService
	save   *saveService
}

func NewService(
	r Repository,
	tr DocumentService,
) *Service {
	return &Service{
		search: newSearchService(r, tr),
		save:   newSaveService(r),
	}
}

func (s *Service) SearchSimilar(
	ctx context.Context,
	query model.DocumentSimilarQuery,
) ([]*model.DocumentSimilarMatch, error) {
	res, err := s.search.SearchSimilar(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	return res, nil
}

func (s *Service) Save(
	ctx context.Context,
	cmd model.DocumentSaveCommand,
) error {
	err := s.save.Save(ctx, cmd)
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}

	return nil
}
