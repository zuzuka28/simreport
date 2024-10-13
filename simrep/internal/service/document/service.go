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
}

func NewService(
	r Repository,
	imgR ImageRepository,
	fileR FileRepository,
) *Service {
	return &Service{
		r:     r,
		imgR:  imgR,
		fileR: fileR,
	}
}

func (s *Service) UploadManyFiles(
	ctx context.Context,
	cmd model.DocumentFileUploadManyCommand,
) error {
	g, gCtx := errgroup.WithContext(ctx)

	for _, item := range cmd.Items {
		g.Go(func() error {
			return s.UploadFile(gCtx, model.DocumentFileUploadCommand{
				Item: item,
			})
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("upload file: %w", err)
	}

	return nil
}

func (s *Service) UploadFile(
	ctx context.Context,
	cmd model.DocumentFileUploadCommand,
) error {
	parsed, err := docxparser.Parse(cmd.Item)
	if err != nil {
		return fmt.Errorf("parse document: %w", err)
	}

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return s.fileR.SaveMany(gCtx, model.MediaFileSaveManyCommand{
			Items: []model.File{parsed.Source},
		})
	})

	g.Go(func() error {
		return s.r.Save(gCtx, model.DocumentSaveCommand{
			Item: mapParsedDocumentFileToDocument(parsed),
		})
	})

	g.Go(func() error {
		return s.imgR.SaveMany(gCtx, model.MediaFileSaveManyCommand{
			Items: parsed.Images,
		})
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("save file resources: %w", err)
	}

	return nil
}

func (s *Service) Fetch(
	ctx context.Context,
	query model.DocumentQuery,
) (model.Document, error) {
	res, err := s.r.Fetch(ctx, query)
	if err != nil {
		return model.Document{}, fmt.Errorf("fetch document: %w", err)
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

func (s *Service) FetchFile(
	ctx context.Context,
	query model.DocumentFileQuery,
) (model.DocumentFile, error) {
	res, err := s.fileR.Fetch(ctx, query)
	if err != nil {
		return model.DocumentFile{}, fmt.Errorf("fetch file: %w", err)
	}

	return res, nil
}

func (s *Service) FetchParsedFile(
	ctx context.Context,
	query model.ParsedDocumentFileQuery,
) (model.ParsedDocumentFile, error) {
	doc, err := s.Fetch(ctx, model.DocumentQuery{
		ID: query.DocumentID,
	})
	if err != nil {
		return model.ParsedDocumentFile{}, fmt.Errorf("fetch document: %w", err)
	}

	g, gCtx := errgroup.WithContext(ctx)

	var main model.DocumentFile

	g.Go(func() error {
		r, err := s.fileR.Fetch(gCtx, model.DocumentFileQuery{
			ID: doc.ID,
		})
		if err != nil {
			return fmt.Errorf("fetch document file: %w", err)
		}

		main = r

		return nil
	})

	var (
		media   []model.MediaFile
		mediaMu sync.Mutex
	)

	for _, id := range doc.ImageIDs {
		g.Go(func() error {
			r, err := s.imgR.Fetch(gCtx, model.MediaFileQuery{
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
		return model.ParsedDocumentFile{}, fmt.Errorf("fetch: %w", err)
	}

	return model.ParsedDocumentFile{
		ID:          doc.ID,
		Source:      main,
		Images:      media,
		TextContent: doc.TextContent,
	}, nil
}
