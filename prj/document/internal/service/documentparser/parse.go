package documentparser

import (
	"bytes"
	"context"
	"document/internal/model"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func (s *Service) Parse(ctx context.Context, item model.File) (model.Document, error) {
	if len(item.Content) == 0 {
		return model.Document{}, errEmptyFile
	}

	doc := model.Document{
		ParentID:    "",
		Name:        item.Name,
		LastUpdated: time.Now(),
		Version:     0,
		GroupID:     nil,
		SourceID:    "",
		TextID:      "",
		ImageIDs:    nil,
		WithContent: true,
		Source:      item,
		Text:        model.File{}, //nolint:exhaustruct
		Images:      nil,
	}

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		text, err := s.tika.ParseText(egCtx, bytes.NewReader(item.Content))
		if err != nil {
			return fmt.Errorf("parse text: %w", err)
		}

		doc.Text = model.File{
			Name:        "",
			Content:     text.Content,
			Sha256:      text.Sha256,
			LastUpdated: doc.LastUpdated,
		}

		doc.TextID = text.Sha256

		return nil
	})

	eg.Go(func() error {
		imgs, err := s.tika.ParseEmbedded(egCtx, bytes.NewReader(item.Content))
		if err != nil {
			return fmt.Errorf("parse images: %w", err)
		}

		images := make([]model.File, len(imgs))
		for i, img := range imgs {
			images[i] = model.File{
				Name:        "",
				Content:     img.Content,
				Sha256:      img.Sha256,
				LastUpdated: doc.LastUpdated,
			}
		}

		doc.Images = images

		imageIDs := make([]string, len(imgs))
		for i, img := range imgs {
			imageIDs[i] = img.Sha256
		}

		doc.ImageIDs = imageIDs

		return nil
	})

	if err := eg.Wait(); err != nil {
		return model.Document{}, fmt.Errorf("parse document: %w", err)
	}

	return doc, nil
}
