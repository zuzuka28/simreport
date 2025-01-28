package shingleindex

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/zuzuka28/simreport/prj/shingleindex/internal/model"

	"github.com/zuzuka28/simreport/lib/minhash"
	"github.com/zuzuka28/simreport/lib/sequencematcher"

	"golang.org/x/sync/errgroup"
)

type searchService struct {
	r  Repository
	tr DocumentService
}

func newSearchService(
	r Repository,
	tr DocumentService,
) *searchService {
	return &searchService{
		r:  r,
		tr: tr,
	}
}

func (s *searchService) SearchSimilar(
	ctx context.Context,
	query model.DocumentSimilarQuery,
) ([]*model.DocumentSimilarMatch, error) {
	nquery := mapDocumentToMinhashSimilarQuery(query)

	res, err := s.r.SearchSimilar(ctx, nquery)
	if err != nil {
		return nil, fmt.Errorf("search similar: %w", err)
	}

	processed, err := s.postprocessing(ctx, query, res)
	if err != nil {
		return nil, fmt.Errorf("postprocessing: %w", err)
	}

	return processed, nil
}

func (s *searchService) postprocessing(
	ctx context.Context,
	query model.DocumentSimilarQuery,
	items []*model.MinhashSimilarMatch,
) ([]*model.DocumentSimilarMatch, error) {
	sources, err := s.fetchSourceDocuments(ctx, items)
	if err != nil {
		return nil, fmt.Errorf("fetch sources: %w", err)
	}

	highlighted := s.highlight(query, sources)

	reranked := s.rerank(query, highlighted)

	result := make([]*model.DocumentSimilarMatch, 0, len(items))

	for _, v := range reranked {
		if v.Rate > 0 {
			result = append(result, v.DocumentSimilarMatch)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Rate > result[j].Rate
	})

	return result, nil
}

func (*searchService) rerank(
	query model.DocumentSimilarQuery,
	items map[string]*documentMatch,
) map[string]*documentMatch {
	sourcesnormalized := normalizeMatches(items)
	sourcesshingled := shingleMatches(sourcesnormalized)

	text := string(query.Item.Text)
	text = normalize(text)
	textshingles := shingle(text, shingleSize)

	for k, v := range sourcesshingled {
		v.Rate = jaccardSimilarity(textshingles, v.shingles)
		sourcesshingled[k] = v
	}

	return items
}

func (*searchService) highlight(
	query model.DocumentSimilarQuery,
	items map[string]*documentMatch,
) map[string]*documentMatch {
	matcher := sequencematcher.NewMatcher[string]()

	text := string(query.Item.Text)
	textwords := strings.Fields(text)

	matcher.SetSeq2(textwords)

	for _, v := range items {
		matcher.SetSeq1(strings.Fields(v.text))

		var highlights []string

		for _, match := range matcher.GetMatchingBlocks() {
			highlight := strings.Join(textwords[match.A:match.A+match.Size], " ")
			if highlight == "" {
				continue
			}

			highlights = append(highlights, highlight)
		}

		v.Highlights = highlights
	}

	return items
}

func (s *searchService) fetchSourceDocuments(
	ctx context.Context,
	matches []*model.MinhashSimilarMatch,
) (map[string]*documentMatch, error) {
	eg, egCtx := errgroup.WithContext(ctx)

	result := make(map[string]*documentMatch, len(matches))
	resultMu := &sync.Mutex{}

	for _, v := range matches {
		eg.Go(func() error {
			f, err := s.tr.Fetch(egCtx, model.DocumentQuery{
				ID: v.DocumentID,
			})
			if err != nil {
				return fmt.Errorf("fetch text: %w", err)
			}

			resultMu.Lock()
			defer resultMu.Unlock()

			result[v.DocumentID] = &documentMatch{
				DocumentSimilarMatch: &model.DocumentSimilarMatch{
					ID:            v.DocumentID,
					Rate:          0,
					Highlights:    nil,
					SimilarImages: nil,
				},
				MinhashSimilarMatch: v,
				text:                string(f.Text),
				shingles:            nil,
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("fetch texts: %w", err)
	}

	return result, nil
}

func normalizeMatches(
	items map[string]*documentMatch,
) map[string]*documentMatch {
	for k, v := range items {
		v.text = normalize(v.text)

		items[k] = v
	}

	return items
}

func shingleMatches(
	items map[string]*documentMatch,
) map[string]*documentMatch {
	for k, v := range items {
		v.shingles = shingle(v.text, shingleSize)
		items[k] = v
	}

	return items
}

func mapDocumentToMinhashSimilarQuery(
	in model.DocumentSimilarQuery,
) model.MinhashSimilarQuery {
	text := string(in.Item.Text)
	text = normalize(text)

	shingles := shingle(text, shingleSize)

	mh := minhash.New(
		permutations,
		hasher,
		seed,
	)

	for shingle := range shingles {
		mh.Push([]byte(shingle))
	}

	return model.MinhashSimilarQuery{
		Minhash: mh,
	}
}
