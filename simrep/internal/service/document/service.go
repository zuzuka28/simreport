package document

import (
	"context"
	"fmt"
	"simrep/internal/model"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

const (
	bucketText  = "texts"
	bucketImage = "images"
)

//nolint:gochecknoglobals
var genID = uuid.NewString

type Opts struct {
	OnSaveAction func(ctx context.Context, cmd model.DocumentSaveCommand) error
}

type Service struct {
	r  Repository
	fr FileRepository
	p  Parser

	onSaveAction func(ctx context.Context, cmd model.DocumentSaveCommand) error
}

func NewService(
	opts Opts,
	r Repository,
	fr FileRepository,
	p Parser,
) *Service {
	return &Service{
		r:            r,
		fr:           fr,
		p:            p,
		onSaveAction: opts.OnSaveAction,
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
		res, err = s.enrichContent(ctx, res, query.Include)
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

func (s *Service) Parse(
	ctx context.Context,
	item model.File,
) (model.Document, error) {
	parsed, err := s.p.Parse(ctx, item)
	if err != nil {
		return model.Document{}, fmt.Errorf("parse document: %w", err)
	}

	return parsed, nil
}

func (s *Service) Save(
	ctx context.Context,
	cmd model.DocumentSaveCommand,
) (*model.Document, error) {
	cmd = s.prepareSaveCommand(cmd)

	g, gCtx := errgroup.WithContext(ctx)

	if cmd.Item.ID != "" {
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
	if cmd.Item.ID == "" {
		cmd.Item.ID = genID()
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
		ID:          doc.ID,
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
