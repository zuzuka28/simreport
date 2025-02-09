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

	processed, err := s.processCandidateMatches(ctx, query, sources)
	if err != nil {
		return nil, fmt.Errorf("refine matches: %w", err)
	}

	if err := s.hr.Save(ctx, model.SimilarityHistorySaveCommand{
		Item: model.SimilarityHistory{
			Date:       now(),
			DocumentID: query.ID,
			ID:         genID(),
			Matches:    processed,
		},
	}); err != nil {
		return nil, fmt.Errorf("save history: %w", err)
	}

	return processed, nil
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

func (s *Service) processCandidateMatches(
	ctx context.Context,
	query model.SimilarityQuery,
	items []*model.SimilarityMatch,
) ([]*model.SimilarityMatch, error) {
	items = deduplicateByID(items)

	docs, err := s.expandSourcesToDocuments(ctx, items)
	if err != nil {
		return nil, fmt.Errorf("expand sources: %w", err)
	}

	sources, err := s.fetchSourceTexts(ctx, docs)
	if err != nil {
		return nil, fmt.Errorf("fetch sources: %w", err)
	}

	highlighted := highlight(string(query.Item.Text.Content), sources)

	reranked := rerank(string(query.Item.Text.Content), highlighted)

	result := make([]*model.SimilarityMatch, 0, len(items))

	for _, m := range reranked {
		if m.Rate <= 0 {
			continue
		}

		for _, d := range m.docs {
			dm := *m
			dm.ID = d.ID()

			result = append(result, dm.SimilarityMatch)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Rate > result[j].Rate
	})

	return result, nil
}

func (s *Service) expandSourcesToDocuments(
	ctx context.Context,
	items []*model.SimilarityMatch,
) (map[string]*match, error) {
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

	result := make(map[string]*match)

	for _, v := range docs {
		m, ok := result[v.SourceID]
		if ok {
			m.docs = append(m.docs, v)
			continue
		}

		result[v.SourceID] = &match{
			SimilarityMatch: srcToMatch[v.SourceID],
			docs:            []model.Document{v},
			textid:          v.TextID,
			text:            "",
			shingles:        nil,
		}
	}

	return result, nil
}

func (s *Service) fetchSourceTexts(
	ctx context.Context,
	matches map[string]*match,
) (map[string]*match, error) {
	eg, egCtx := errgroup.WithContext(ctx)

	for _, v := range matches {
		eg.Go(func() error {
			f, err := s.fs.Fetch(egCtx, model.FileQuery{
				Bucket: bucketTexts,
				ID:     v.textid,
			})
			if err != nil {
				return fmt.Errorf("fetch text: %w", err)
			}

			matches[v.ID].text = string(f.Content)

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("fetch texts: %w", err)
	}

	return matches, nil
}
