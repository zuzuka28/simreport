// anysave - saves everything that is sent to the simrep service.
package anysave

import (
	"context"
	"fmt"
	"simrep/internal/model"
)

const (
	bucketAnysave = "anysave"
)

type Opts struct {
	OnSaveAction func(ctx context.Context, cmd model.FileSaveCommand) error
}

type Service struct {
	r            Repository
	onSaveAction func(ctx context.Context, cmd model.FileSaveCommand) error
}

func NewService(
	opts Opts,
	r Repository,
) *Service {
	return &Service{
		r:            r,
		onSaveAction: opts.OnSaveAction,
	}
}

func (s *Service) Fetch(
	ctx context.Context,
	query model.FileQuery,
) (model.File, error) {
	if query.Bucket == "" {
		query.Bucket = bucketAnysave
	}

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
	if cmd.Bucket == "" {
		cmd.Bucket = bucketAnysave
	}

	_, err := s.r.Fetch(ctx, model.FileQuery{
		Bucket: cmd.Bucket,
		ID:     cmd.Item.Sha256,
	})
	if err == nil { // == nil
		return nil
	}

	if err := s.r.Save(ctx, cmd); err != nil {
		return fmt.Errorf("save file: %w", err)
	}

	if s.onSaveAction != nil {
		if err := s.onSaveAction(ctx, cmd); err != nil {
			return fmt.Errorf("on save action: %w", err)
		}
	}

	return nil
}
