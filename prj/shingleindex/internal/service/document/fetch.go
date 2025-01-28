package document

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/shingleindex/internal/model"
)

func (s *Service) Fetch(
	ctx context.Context,
	query model.DocumentQuery,
) (model.Document, error) {
	res, err := s.r.Fetch(ctx, query)
	if err != nil {
		return model.Document{}, fmt.Errorf("fetch document: %w", err)
	}

	return res, nil
}
