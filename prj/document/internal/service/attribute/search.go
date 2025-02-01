package attribute

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/document/internal/model"
)

func (s *Service) Search(
	ctx context.Context,
	query model.AttributeQuery,
) ([]model.Attribute, error) {
	res, err := s.r.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("search attribute: %w", err)
	}

	return res, nil
}
