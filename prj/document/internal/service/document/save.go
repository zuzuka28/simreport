package document

import (
	"context"
	"document/internal/model"
	"fmt"

	"golang.org/x/sync/errgroup"
)

func (s *Service) Save(
	ctx context.Context,
	cmd model.DocumentSaveCommand,
) (*model.Document, error) {
	cmd = s.prepareSaveCommand(cmd)

	g, gCtx := errgroup.WithContext(ctx)

	if cmd.Item.ParentID != "" {
		g.Go(func() error {
			return s.r.Save(gCtx, model.DocumentSaveCommand{
				Item: mapDocumentWithContentToDocument(cmd.Item),
			})
		})
	}

	if cmd.Item.Text.Sha256 != "" {
		g.Go(func() error {
			return s.fr.Save(gCtx, model.FileSaveCommand{
				Bucket: bucketText,
				Item:   cmd.Item.Text,
			})
		})
	}

	if cmd.Item.Source.Sha256 != "" {
		g.Go(func() error {
			return s.fr.Save(gCtx, model.FileSaveCommand{
				Bucket: "",
				Item:   cmd.Item.Source,
			})
		})
	}

	for _, img := range cmd.Item.Images {
		if img.Sha256 != "" {
			g.Go(func() error {
				return s.fr.Save(gCtx, model.FileSaveCommand{
					Bucket: bucketImage,
					Item:   img,
				})
			})
		}
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("save file resources: %w", err)
	}

	if s.onSaveAction != nil {
		if err := s.onSaveAction(ctx, cmd); err != nil {
			return nil, fmt.Errorf("on save action: %w", err)
		}
	}

	return &cmd.Item, nil
}

func (*Service) prepareSaveCommand(
	cmd model.DocumentSaveCommand,
) model.DocumentSaveCommand {
	if cmd.Item.ParentID == "" {
		cmd.Item.ParentID = genID()
	}

	if cmd.Item.Source.Sha256 != "" {
		cmd.Item.SourceID = cmd.Item.Source.Sha256
	}

	if cmd.Item.Text.Sha256 != "" {
		cmd.Item.TextID = cmd.Item.Text.Sha256
	}

	if len(cmd.Item.Images) > 0 {
		imgIDs := make([]string, 0, len(cmd.Item.Images))
		for _, v := range cmd.Item.Images {
			imgIDs = append(imgIDs, v.Sha256)
		}

		cmd.Item.ImageIDs = imgIDs
	}

	return cmd
}
