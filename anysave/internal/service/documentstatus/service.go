package documentstatus

import (
	"context"
	"fmt"
	"anysave/internal/model"
)

type Service struct {
	r Repository
}

func NewService(r Repository) *Service {
	return &Service{
		r: r,
	}
}

func (s *Service) Update(
	ctx context.Context,
	cmd model.DocumentStatusUpdateCommand,
) error {
	if err := s.r.Update(ctx, cmd); err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	return nil
}

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
