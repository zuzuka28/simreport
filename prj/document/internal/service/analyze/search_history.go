package analyze

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/document/internal/model"
)

func (s *Service) SearchHistory(
	ctx context.Context,
	query model.SimilarityHistoryQuery,
) (*model.SimilarityHistoryList, error) {
	res, err := s.hr.Fetch(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("fetch history: %w", err)
	}

	return res, nil
}
