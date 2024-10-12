package docprepsrv

import (
	"context"
	"fmt"
	"simrep/internal/model"
)

type Service struct {
	r Repository
}

func NewService(r Repository) *Service {
	return &Service{
		r: r,
	}
}

func (s *Service) PreprocessRawDocument(
	ctx context.Context,
	doc []byte,
) (*model.DocumentFull, error) {
	res, err := s.r.PreprocessRawDocument(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("do preprocess raw: %w", err)
	}

	return res, nil
}
