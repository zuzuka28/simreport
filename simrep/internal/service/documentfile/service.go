package documentfile

import (
	"context"
	"fmt"
	"simrep/internal/model"
)

type Service struct {
	r Repository
	n     Notify
}

func NewService(
	r Repository,
	n Notify,
) *Service {
	return &Service{
		r: r,
		n:     n,
	}
}

func (s *Service) Fetch(
	ctx context.Context,
	query model.FileQuery,
) (model.File, error) {
	res, err := s.r.Fetch(ctx, query)
	if err != nil {
		return model.File{}, fmt.Errorf("fetch document file: %w", err)
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
		return fmt.Errorf("save document file: %w", err)
	}

	if err := s.n.Notify(
		ctx,
		cmd.Item.Sha256,
		model.NotifyActionFileSaved,
		nil,
	); err != nil {
		return fmt.Errorf("notify: %w", err)
	}

	return nil
}
