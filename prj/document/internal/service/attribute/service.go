package attribute

import (
	"context"
	"document/internal/model"
	"fmt"
)

type Service struct {
	r Repository
}

func NewService(
	r Repository,
) *Service {
	return &Service{
		r: r,
	}
}

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
