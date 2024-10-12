package document

import (
	"context"
	"fmt"
	"simrep/internal/model"

	"golang.org/x/sync/errgroup"
)

type Service struct {
	parser FileParser
	r      Repository
	imgR   ImageRepository
	fileR  FileRepository
}

func NewService(
	parser FileParser,
	r Repository,
	imgR ImageRepository,
	fileR FileRepository,
) *Service {
	return &Service{
		parser: parser,
		r:      r,
		imgR:   imgR,
		fileR:  fileR,
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
	parsed, err := s.parser.Parse(ctx, cmd.Item)
	if err != nil {
		return fmt.Errorf("parse document: %w", err)
	}

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return s.fileR.SaveMany(gCtx, model.MediaFileSaveManyCommand{
			Items: []model.File{cmd.Item},
		})
	})

	g.Go(func() error {
		return s.r.SaveParsed(gCtx, model.ParsedDocumentSaveCommand{
			Item: mapParsedDocumentFileToParsedDocument(parsed),
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
