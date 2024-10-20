package document

import (
	"context"
	"fmt"
	"simrep/internal/model"
	"simrep/internal/service/document/docxparser"
	"sync"

	"golang.org/x/sync/errgroup"
)

type Service struct {
	r     Repository
	imgR  ImageRepository
	fileR FileRepository
	n     Notify
}

func NewService(
	r Repository,
	imgR ImageRepository,
	fileR FileRepository,
	n Notify,
) *Service {
	return &Service{
		r:     r,
		imgR:  imgR,
		fileR: fileR,
		n:     n,
	}
}

func (s *Service) Fetch(
	ctx context.Context,
	query model.DocumentQuery,
) (model.Document, error) {
	res, err := s.r.Fetch(ctx, query)
	if err != nil {
		return model.Document{}, fmt.Errorf("fetch document: %w", err)
	}

	if query.WithContent {
		res, err = s.enrichContent(ctx, res)
		if err != nil {
			return model.Document{}, fmt.Errorf("enrich document: %w", err)
		}
	}

	return res, nil
}

func (s *Service) Search(
	ctx context.Context,
	query model.DocumentSearchQuery,
) ([]model.Document, error) {
	res, err := s.r.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("search documents: %w", err)
	}

	return res, nil
}

func (*Service) Parse(
	_ context.Context,
	item model.File,
) (model.Document, error) {
	parsed, err := docxparser.Parse(item)
	if err != nil {
		return model.Document{}, fmt.Errorf("parse document: %w", err)
	}

	return parsed, nil
}

func (s *Service) Save(
	ctx context.Context,
	cmd model.DocumentSaveCommand,
) error {
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return s.r.Save(gCtx, model.DocumentSaveCommand{
			Item: mapDocumentWithContentToDocument(cmd.Item),
		})
	})

	g.Go(func() error {
		return s.fileR.Save(gCtx, model.FileSaveCommand{
			Item: cmd.Item.Source,
		})
	})

	for _, img := range cmd.Item.Images {
		g.Go(func() error {
			return s.imgR.Save(gCtx, model.FileSaveCommand{
				Item: img,
			})
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("save file resources: %w", err)
	}

	if err := s.n.Notify(
		ctx,
		cmd.Item.ID,
		model.NotifyActionDocumentSaved,
		nil,
	); err != nil {
		return fmt.Errorf("notify: %w", err)
	}

	return nil
}

func (s *Service) enrichContent(
	ctx context.Context,
	doc model.Document,
) (model.Document, error) {
	g, gCtx := errgroup.WithContext(ctx)

	var main model.File

	g.Go(func() error {
		r, err := s.fileR.Fetch(gCtx, model.FileQuery{
			ID: doc.ID,
		})
		if err != nil {
			return fmt.Errorf("fetch document file: %w", err)
		}

		main = r

		return nil
	})

	var (
		media   []model.File
		mediaMu sync.Mutex
	)

	for _, id := range doc.ImageIDs {
		g.Go(func() error {
			r, err := s.imgR.Fetch(gCtx, model.FileQuery{
				ID: id,
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

	if err := g.Wait(); err != nil {
		return model.Document{}, fmt.Errorf("fetch: %w", err)
	}

	return model.Document{
		ID:          doc.ID,
		Name:        doc.Name,
		ImageIDs:    doc.ImageIDs,
		TextContent: doc.TextContent,
		LastUpdated: doc.LastUpdated,
		WithContent: true,
		Source:      main,
		Images:      media,
	}, nil
}
