package analyzehistory

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/zuzuka28/simreport/prj/similarity/internal/model"

	"github.com/zuzuka28/simreport/lib/elasticutil"
)

func (r *Repository) Fetch(
	ctx context.Context,
	query model.SimilarityHistoryQuery,
) (*model.SimilarityHistoryList, error) {
	const op = "fetch"

	q, err := buildSearchQuery(query)
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
		r.m.IncAnalyzeHistoryRepositoryRequests(op, metricsError, time.Since(t).Seconds())
		return nil, fmt.Errorf("search history: %w", err)
	}

	defer esRes.Body.Close()

	r.m.IncAnalyzeHistoryRepositoryRequests(op, esRes.Status(), time.Since(t).Seconds())

	if err := elasticutil.IsErr(esRes); err != nil {
		return nil, fmt.Errorf("search error: %s: %w", esRes.Status(), mapErrorToModel(err))
	}

	raw, err := elasticutil.ParseSearchResponse(esRes.Body)
	if err != nil {
		return nil, fmt.Errorf("parse search response: %w", err)
	}

	res, err := mapSearchResponseToHistoryList(raw)
	if err != nil {
		return nil, fmt.Errorf("map search response to history list: %w", err)
	}

	return res, nil
}
