package shingleindex

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/shingleindex/internal/model"
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

	return mapMinhashMatchesToDocumentMatches(res), nil
}
