package shingleindex

import (
	"context"
	"fmt"
	"regexp"
	"simrep/internal/model"
	"sort"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

const shingleSize = 4

type Service struct {
	r  Repository
	tr DocumentService
}

func NewService(
	r Repository,
	tr DocumentService,
) *Service {
	return &Service{
		r:  r,
		tr: tr,
	}
}

func (s *Service) SearchSimilar(
	ctx context.Context,
	query model.DocumentSimilarQuery,
) ([]model.DocumentSimilarMatch, error) {
	res, err := s.r.SearchSimilar(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("search similar: %w", err)
	}

	res, err = s.rerank(ctx, query, res)
	if err != nil {
		return nil, fmt.Errorf("rerank similar: %w", err)
	}

	var nonzero []model.DocumentSimilarMatch

	for _, v := range res {
		if v.Rate > 0 {
			nonzero = append(nonzero, v)
		}
	}

	return nonzero, nil
}

func (s *Service) rerank(
	ctx context.Context,
	query model.DocumentSimilarQuery,
	items []model.DocumentSimilarMatch,
) ([]model.DocumentSimilarMatch, error) {
	sources, err := s.fetchSourceDocuments(ctx, items)
	if err != nil {
		return nil, fmt.Errorf("fetch sources: %w", err)
	}

	sourcesnormalized := normalizeMatches(sources)
	sourcesshingled := shingleMatches(sourcesnormalized)

	text := string(query.Item.Text.Content)
	text = normalize(text)
	textshingles := shingle(text, shingleSize)

	for k, v := range sourcesshingled {
		v.Rate = jaccardSimilarity(textshingles, v.shingles)
		sourcesshingled[k] = v
	}

	reranked := make([]model.DocumentSimilarMatch, 0, len(items))
	for _, v := range sourcesshingled {
		reranked = append(reranked, v.DocumentSimilarMatch)
	}

	sort.Slice(reranked, func(i, j int) bool {
		return reranked[i].Rate > reranked[j].Rate
	})

	return reranked, nil
}

func (s *Service) fetchSourceDocuments(
	ctx context.Context,
	matches []model.DocumentSimilarMatch,
) (map[string]*documentMatch, error) {
	eg, egCtx := errgroup.WithContext(ctx)

	result := make(map[string]*documentMatch, len(matches))
	resultMu := &sync.Mutex{}

	for _, v := range matches {
		eg.Go(func() error {
			f, err := s.tr.Fetch(egCtx, model.DocumentQuery{
				ID:          v.ID,
				WithContent: true,
				Include:     []model.DocumentQueryInclude{model.DocumentQueryIncludeText},
			})
			if err != nil {
				return fmt.Errorf("fetch text: %w", err)
			}

			resultMu.Lock()
			defer resultMu.Unlock()

			result[v.ID] = &documentMatch{
				DocumentSimilarMatch: v,
				text:                 string(f.Text.Content),
				shingles:             nil,
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

var (
	reNormalize1 = regexp.MustCompile(`[^a-zA-Zа-яА-ЯёЁ0-9_ ]`)
	reNormalize2 = regexp.MustCompile(`\s+`)
)

func normalize(text string) string {
	text = reNormalize1.ReplaceAllString(text, "")
	text = reNormalize2.ReplaceAllString(text, " ")
	text = strings.ToLower(text)

	return text
}

func shingle(text string, shingleSize int) map[string]struct{} {
	words := strings.Fields(text)

	var shingles []string

	for i := 0; i <= len(words)-shingleSize; i++ {
		shingle := strings.Join(words[i:i+shingleSize], " ")
		shingles = append(shingles, shingle)
	}

	shingleSet := make(map[string]struct{})
	for _, shingle := range shingles {
		shingleSet[shingle] = struct{}{}
	}

	return shingleSet
}

func jaccardSimilarity(set1, set2 map[string]struct{}) float64 {
	intersection := 0

	for item := range set1 {
		if _, exists := set2[item]; exists {
			intersection++
		}
	}

	union := len(set1) + len(set2) - intersection

	if union == 0 {
		return 0
	}

	return float64(intersection) / float64(union)
}
