package document

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"
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

func (s *Service) Save(
	ctx context.Context,
	cmd model.DocumentSaveCommand,
) (*model.DocumentSaveResult, error) {
	res, err := s.r.Save(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("save document: %w", err)
	}

	return res, nil
}
