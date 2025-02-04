package vectorizer

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/semanticindex/internal/model"
)

func (s *Service) TextToVector(
	ctx context.Context,
	params model.VectorizeTextParams,
) (model.Vector, error) {
	res, err := s.r.TextToVector(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("vectorize: %w", err)
	}

	return res, nil
}
