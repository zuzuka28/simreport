package similarity

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/zuzuka28/simreport/prj/similarity/internal/model"

	"golang.org/x/sync/errgroup"
)

const bucketTexts = "texts"

func (s *Service) SearchSimilar(
	ctx context.Context,
	query model.SimilarityQuery,
) ([]*model.SimilarityMatch, error) {
	query, err := s.prepareSearchQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("prepare search query: %w", err)
	}

	sources, err := s.searchSources(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("search sources: %w", err)
	}

	sources, err = s.processCandidateMatches(ctx, query, sources)
	if err != nil {
		return nil, fmt.Errorf("refine matches: %w", err)
	}

	res, err := s.expandSourcesToDocuments(ctx, sources)
	if err != nil {
		return nil, fmt.Errorf("expand matches: %w", err)
	}

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

func (s *Service) prepareSearchQuery(
	ctx context.Context,
	query model.SimilarityQuery,
) (model.SimilarityQuery, error) {
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
		return model.SimilarityQuery{}, fmt.Errorf("enrich query with document:%w", err)
	}

	query.Item = doc

	return query, nil
}

func (s *Service) searchSources(
	ctx context.Context,
	query model.SimilarityQuery,
) ([]*model.SimilarityMatch, error) {
	eg, egCtx := errgroup.WithContext(ctx)

	var (
		res   []*model.SimilarityMatch
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

	return res, nil
}

func (s *Service) expandSourcesToDocuments(
	ctx context.Context,
	items []*model.SimilarityMatch,
) ([]*model.SimilarityMatch, error) {
	srcToMatch := make(map[string]*model.SimilarityMatch)

	for _, v := range items {
		srcToMatch[v.ID] = v
	}

	docs, err := s.ds.Search(ctx, model.DocumentSearchQuery{ //nolint:exhaustruct
		SourceID: extractSourceIDs(items),
	})
	if err != nil {
		return nil, fmt.Errorf("retrieve document: %w", err)
	}

	result := make([]*model.SimilarityMatch, 0, len(docs))

	for _, v := range docs {
		m := *srcToMatch[v.SourceID]
		m.ID = v.ID()

		result = append(result, &m)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Rate > result[j].Rate
	})

	return result, nil
}

func (s *Service) processCandidateMatches(
	ctx context.Context,
	query model.SimilarityQuery,
	items []*model.SimilarityMatch,
) ([]*model.SimilarityMatch, error) {
	items = deduplicateByID(items)

	sources, err := s.fetchSourceTexts(ctx, items)
	if err != nil {
		return nil, fmt.Errorf("fetch sources: %w", err)
	}

	highlighted := highlight(string(query.Item.Text.Content), sources)

	reranked := rerank(string(query.Item.Text.Content), highlighted)

	result := make([]*model.SimilarityMatch, 0, len(items))

	for _, v := range reranked {
		if v.Rate > 0 {
			result = append(result, v.SimilarityMatch)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Rate > result[j].Rate
	})

	return result, nil
}

func (s *Service) fetchSourceTexts(
	ctx context.Context,
	matches []*model.SimilarityMatch,
) (map[string]*match, error) {
	eg, egCtx := errgroup.WithContext(ctx)

	result := make(map[string]*match, len(matches))
	resultMu := &sync.Mutex{}

	for _, v := range matches {
		eg.Go(func() error {
			f, err := s.fs.Fetch(egCtx, model.FileQuery{
				Bucket: bucketTexts,
				ID:     v.ID,
			})
			if err != nil {
				return fmt.Errorf("fetch text: %w", err)
			}

			resultMu.Lock()
			defer resultMu.Unlock()

			result[v.ID] = &match{
				SimilarityMatch: &model.SimilarityMatch{
					ID:            v.ID,
					Rate:          0,
					Highlights:    nil,
					SimilarImages: nil,
				},
				text:     string(f.Content),
				shingles: nil,
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("fetch texts: %w", err)
	}

	return result, nil
}
