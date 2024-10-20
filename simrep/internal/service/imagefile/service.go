package imagefile

import (
	"context"
	"fmt"
	"simrep/internal/model"
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
	query model.FileQuery,
) (model.File, error) {
	res, err := s.r.Fetch(ctx, query)
	if err != nil {
		return model.File{}, fmt.Errorf("fetch image file: %w", err)
	}

	return res, nil
}

func (s *Service) Save(
	ctx context.Context,
	cmd model.FileSaveCommand,
) error {
	_, err := s.r.Fetch(ctx, model.FileQuery{
		ID: cmd.Item.Sha256,
	})
	if err == nil { // == nil
		return nil
	}

	if err := s.r.Save(ctx, cmd); err != nil {
		return fmt.Errorf("save image file: %w", err)
	}

	return nil
}
