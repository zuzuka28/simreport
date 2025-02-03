package similarity

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"
)

func (s *Repository) Search(
	ctx context.Context,
	query model.SimilarityQuery,
) ([]*model.SimilarityMatch, error) {
	q := mapDocumentSearchQueryToPb(query)

	resp, err := s.cli.SearchSimilar(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	if err := isErr(resp.GetError()); err != nil {
		return nil, err
	}

	return mapSimilarityMatchesToModel(resp), nil
}
