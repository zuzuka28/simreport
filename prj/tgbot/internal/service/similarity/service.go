package similarity

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"
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

func (s *Service) Search(
	ctx context.Context,
	query model.SimilarityQuery,
) ([]*model.SimilarityMatch, error) {
	res, err := s.r.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("search similar: %w", err)
	}

	return res, nil
}
