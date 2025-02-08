package similarity

import (
	"context"
	"fmt"
	"time"

	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"
)

func (s *Repository) Search(
	ctx context.Context,
	query model.SimilarityQuery,
) ([]*model.SimilarityMatch, error) {
	const op = "save"

	q := mapDocumentSearchQueryToPb(query)

	t := time.Now()

	resp, err := s.cli.SearchSimilar(ctx, q)
	if err != nil {
		s.m.IncSimilarityRepositoryRequests(op, mapErrorToStatus(err), time.Since(t).Seconds())
		return nil, fmt.Errorf("do request: %w", mapErrorToModel(err))
	}

	s.m.IncSimilarityRepositoryRequests(op, "200", time.Since(t).Seconds())

	return mapSimilarityMatchesToModel(resp), nil
}
