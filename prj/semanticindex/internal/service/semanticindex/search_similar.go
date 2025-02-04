package semanticindex

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/semanticindex/internal/model"
)

func (s *Service) SearchSimilar(
	ctx context.Context,
	query model.DocumentSimilarQuery,
) ([]model.DocumentSimilarMatch, error) {
	if query.Item.Vector == nil {
		vec, err := s.vs.TextToVector(ctx, model.VectorizeTextParams{
			Text: string(query.Item.Text),
		})
		if err != nil {
			return nil, fmt.Errorf("vectorize document: %w", err)
		}

		query.Item.Vector = vec
	}

	res, err := s.r.SearchSimilar(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("search fulltext similar: %w", err)
	}

	return res, nil
}
