package attribute

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/document/internal/model"
)

func (s *Service) Fetch(
	ctx context.Context,
	query model.AttributeQuery,
) ([]model.Attribute, error) {
	res, err := s.r.Fetch(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("fetch attribute: %w", err)
	}

	return res, nil
}
