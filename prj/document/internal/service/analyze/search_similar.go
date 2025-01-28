package analyze

import (
	"context"
	"document/internal/model"
	"fmt"
	"sort"
	"sync"

	"golang.org/x/sync/errgroup"
)

func (s *Service) SearchSimilar(
	ctx context.Context,
	query model.DocumentSimilarQuery,
) ([]*model.DocumentSimilarMatch, error) {
	doc, err := s.ds.Fetch(ctx, model.DocumentQuery{
		ID:          query.ID,
		WithContent: true,
		Include: []model.DocumentQueryInclude{
			model.DocumentQueryIncludeSource,
			model.DocumentQueryIncludeText,
			model.DocumentQueryIncludeImages,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("enrich query with document:%w", err)
	}

	query.Item = doc

	eg, egCtx := errgroup.WithContext(ctx)

	var (
		res   []*model.DocumentSimilarMatch
		resMu sync.Mutex
	)

	eg.Go(func() error {
		r, err := s.shingleis.SearchSimilar(egCtx, query)
		if err != nil {
			return fmt.Errorf("shingle similar: %w", err)
		}

		resMu.Lock()
		defer resMu.Unlock()

		res = append(res, r...)

		return nil
	})

	eg.Go(func() error {
		r, err := s.fulltextis.SearchSimilar(egCtx, query)
		if err != nil {
			return fmt.Errorf("fulltext similar: %w", err)
		}

		resMu.Lock()
		defer resMu.Unlock()

		res = append(res, r...)

		return nil
	})

	eg.Go(func() error {
		r, err := s.semanticis.SearchSimilar(egCtx, query)
		if err != nil {
			return fmt.Errorf("semantic similar: %w", err)
		}

		resMu.Lock()
		defer resMu.Unlock()

		res = append(res, r...)

		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("search similar: %w", err)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Rate > res[j].Rate
	})

	if err := s.hr.Save(ctx, model.SimilarityHistorySaveCommand{
		Item: model.SimilarityHistory{
			Date:       now(),
			DocumentID: query.ID,
			ID:         genID(),
			Matches:    res,
		},
	}); err != nil {
		return nil, fmt.Errorf("save history: %w", err)
	}

	return res, nil
}
