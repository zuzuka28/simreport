package analyze

import (
	"fmt"
	"simrep/internal/model"
	"sync"

	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	vs VectorizerService
	r  Repository
	n  Notify
}

func NewService(
	r Repository,
	vs VectorizerService,
	n Notify,
) *Service {
	return &Service{
		vs: vs,
		r:  r,
		n:  n,
	}
}

func (s *Service) Fetch(
	ctx context.Context,
	query model.AnalyzedDocumentQuery,
) (model.AnalyzedDocument, error) {
	res, err := s.r.Fetch(ctx, query)
	if err != nil {
		return model.AnalyzedDocument{}, fmt.Errorf("fetch analyzed document: %w", err)
	}

	return res, nil
}

func (s *Service) SearchSimilar(
	ctx context.Context,
	query model.AnalyzedDocumentSimilarQuery,
) ([]model.AnalyzedDocumentMatch, error) {
	res, err := s.r.SearchSimilar(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("search similar: %w", err)
	}

	return res, nil
}

func (s *Service) Save(
	ctx context.Context,
	cmd model.AnalyzedDocumentSaveCommand,
) error {
	if err := s.r.Save(ctx, cmd); err != nil {
		return fmt.Errorf("save analyzed document: %w", err)
	}

	if err := s.n.Notify(
		ctx,
		cmd.Item.ID,
		model.NotifyActionDocumentAnalyzed,
		nil,
	); err != nil {
		return fmt.Errorf("notify: %w", err)
	}

	return nil
}

func (s *Service) Analyze(
	ctx context.Context,
	item model.Document,
) (model.AnalyzedDocument, error) {
	g, gCtx := errgroup.WithContext(ctx)

	var textVector model.Vector

	g.Go(func() error {
		r, err := s.vs.TextToVector(gCtx, model.VectorizeTextParams{
			Text: item.TextID,
		})
		if err != nil {
			return fmt.Errorf("vectorize text: %w", err)
		}

		textVector = r

		return nil
	})

	var (
		imgs   = make([]model.AnalyzedImage, 0, len(item.Images))
		imgsMu sync.Mutex
	)

	for _, img := range item.Images {
		g.Go(func() error {
			r, err := s.AnalyzeImage(gCtx, img)
			if err != nil {
				return fmt.Errorf("analyze image: %w", err)
			}

			imgsMu.Lock()
			imgs = append(imgs, r)
			imgsMu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return model.AnalyzedDocument{}, fmt.Errorf("analyze document: %w", err)
	}

	return model.AnalyzedDocument{
		ID:         item.ID,
		Text:       item.TextID,
		TextVector: textVector,
		Images:     imgs,
	}, nil
}

func (s *Service) AnalyzeImage(
	ctx context.Context,
	item model.File,
) (model.AnalyzedImage, error) {
	g, gCtx := errgroup.WithContext(ctx)

	var vector model.Vector

	g.Go(func() error {
		r, err := s.vs.ImageToVector(gCtx, model.VectorizeImageParams{
			Image: item,
		})
		if err != nil {
			return fmt.Errorf("vectorize image: %w", err)
		}

		vector = r

		return nil
	})

	var hash model.HashImage

	g.Go(func() error {
		r, err := s.vs.ImageToHashes(gCtx, model.HashImageParams{
			Image: item,
		})
		if err != nil {
			return fmt.Errorf("hash image: %w", err)
		}

		hash = r

		return nil
	})

	if err := g.Wait(); err != nil {
		return model.AnalyzedImage{}, fmt.Errorf("analyze image: %w", err)
	}

	return model.AnalyzedImage{
		ID:        item.Sha256,
		Vector:    vector,
		HashImage: hash,
	}, nil
}
