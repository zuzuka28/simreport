package documentstatus

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/document/internal/model"
)

func (s *Service) Fetch(
	ctx context.Context,
	query model.DocumentStatusQuery,
) ([]*model.DocumentStatus, error) {
	res, err := s.r.Fetch(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("fetch status: %w", err)
	}

	return res, nil
}
