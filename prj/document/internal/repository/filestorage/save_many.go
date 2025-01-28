package filestorage

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/document/internal/model"

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
