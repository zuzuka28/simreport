package fulltextindex

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/zuzuka28/simreport/lib/elasticutil"
	"github.com/zuzuka28/simreport/prj/fulltextindex/internal/model"
)

func (r *Repository) SearchSimilar(
	ctx context.Context,
	query model.DocumentSimilarQuery,
) ([]model.DocumentSimilarMatch, error) {
	const op = "searchSimilar"

	q, err := buildSearchQuery(r.index, query)
	if err != nil {
		return nil, fmt.Errorf("build search query: %w", err)
	}

	t := time.Now()

	esRes, err := r.cli.Search(
		r.cli.Search.WithContext(ctx),
		r.cli.Search.WithIndex(r.index),
		r.cli.Search.WithBody(bytes.NewReader(q)),
	)
	if err != nil {
		r.m.IncFulltextIndexRequests(op, metricsError, time.Since(t).Seconds())
		return nil, fmt.Errorf("search documents: %w", err)
	}

	defer esRes.Body.Close()

	r.m.IncFulltextIndexRequests(op, strconv.Itoa(esRes.StatusCode), time.Since(t).Seconds())

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
