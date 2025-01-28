package documentstatus

import (
	"context"
	"document/internal/model"
	"fmt"
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
