package document

import (
	"context"
	"fmt"
	"sync"

	"github.com/zuzuka28/simreport/prj/document/internal/model"

	"golang.org/x/sync/errgroup"
)

func (s *Service) Fetch(
	ctx context.Context,
	query model.DocumentQuery,
) (model.Document, error) {
	res, err := s.r.Fetch(ctx, query)
	if err != nil {
		return model.Document{}, fmt.Errorf("fetch document: %w", err)
	}

	if query.WithContent {
		res, err = s.enrichContent(ctx, res, query.Include)
		if err != nil {
			return model.Document{}, fmt.Errorf("enrich document: %w", err)
		}
	}

	return res, nil
}

//nolint:revive
func (s *Service) enrichContent(
	ctx context.Context,
	doc model.Document,
	include []model.DocumentQueryInclude,
) (model.Document, error) {
	includemap := make(map[model.DocumentQueryInclude]bool)
	for _, v := range include {
		includemap[v] = true
	}

	eg, egCtx := errgroup.WithContext(ctx)

	var main model.File

	if doc.SourceID != "" && includemap[model.DocumentQueryIncludeSource] {
		eg.Go(func() error {
			r, err := s.fr.Fetch(egCtx, model.FileQuery{
				Bucket: "",
				ID:     doc.SourceID,
			})
			if err != nil {
				return fmt.Errorf("fetch document file: %w", err)
			}

			main = r

			return nil
		})
	}

	var text model.File

	if doc.TextID != "" && includemap[model.DocumentQueryIncludeText] {
		eg.Go(func() error {
			r, err := s.fr.Fetch(egCtx, model.FileQuery{
				Bucket: bucketText,
				ID:     doc.TextID,
			})
			if err != nil {
				return fmt.Errorf("fetch text file: %w", err)
			}

			text = r

			return nil
		})
	}

	var (
		media   []model.File
		mediaMu sync.Mutex
	)

	if includemap[model.DocumentQueryIncludeImages] {
		for _, id := range doc.ImageIDs {
			eg.Go(func() error {
				r, err := s.fr.Fetch(egCtx, model.FileQuery{
					Bucket: bucketImage,
					ID:     id,
				})
				if err != nil {
					return fmt.Errorf("fetch image file: %w", err)
				}

				mediaMu.Lock()
				media = append(media, r)
				mediaMu.Unlock()

				return nil
			})
		}
	}

	if err := eg.Wait(); err != nil {
		return model.Document{}, fmt.Errorf("fetch: %w", err)
	}

	return model.Document{
		ParentID:    doc.ParentID,
		Name:        doc.Name,
		LastUpdated: doc.LastUpdated,
		Version:     doc.Version,
		GroupID:     doc.GroupID,
		SourceID:    doc.SourceID,
		TextID:      doc.TextID,
		ImageIDs:    doc.ImageIDs,
		WithContent: true,
		Source:      main,
		Text:        text,
		Images:      media,
	}, nil
}
