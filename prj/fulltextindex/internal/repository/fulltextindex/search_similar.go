package fulltextindex

import (
	"bytes"
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/lib/elasticutil"
	"github.com/zuzuka28/simreport/prj/fulltextindex/internal/model"
)

func (r *Repository) SearchSimilar(
	ctx context.Context,
	query model.DocumentSimilarQuery,
) ([]model.DocumentSimilarMatch, error) {
	q, err := buildSearchQuery(r.index, query)
	if err != nil {
		return nil, fmt.Errorf("build search query: %w", err)
	}

	esRes, err := r.cli.Search(
		r.cli.Search.WithContext(ctx),
		r.cli.Search.WithIndex(r.index),
		r.cli.Search.WithBody(bytes.NewReader(q)),
	)
	if err != nil {
		return nil, fmt.Errorf("search documents: %w", err)
	}

	defer esRes.Body.Close()

	if err := elasticutil.IsErr(esRes); err != nil {
		return nil, fmt.Errorf("search error: %s: %w", esRes.Status(), err)
	}

	raw, err := elasticutil.ParseSearchResponse(esRes.Body)
	if err != nil {
		return nil, fmt.Errorf("parse search response: %w", err)
	}

	res, err := mapSearchResponseToMatches(raw)
	if err != nil {
		return nil, fmt.Errorf("map search response to documents: %w", err)
	}

	return res, nil
}
