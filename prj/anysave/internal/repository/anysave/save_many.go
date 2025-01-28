package anysave

import (
	"anysave/internal/model"
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

func (r *Repository) SaveMany(ctx context.Context, cmd model.FileSaveManyCommand) error {
	g, gCtx := errgroup.WithContext(ctx)

	for _, item := range cmd.Items {
		g.Go(func() error {
			return r.Save(gCtx, model.FileSaveCommand{
				Bucket: cmd.Bucket,
				Item:   item,
			})
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("save file: %w", err)
	}

	return nil
}
