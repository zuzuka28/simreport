package document

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/zuzuka28/simreport/prj/document/internal/model"

	"github.com/zuzuka28/simreport/lib/elasticutil"
)

func (r *Repository) Search(
	ctx context.Context,
	query model.DocumentSearchQuery,
) ([]model.Document, error) {
	const op = "search"

	t := time.Now()

	q, err := buildSearchQuery(query)
	if err != nil {
		return nil, fmt.Errorf("build search query: %w", err)
	}

	esRes, err := r.cli.Search(
		r.cli.Search.WithContext(ctx),
		r.cli.Search.WithIndex(r.index),
		r.cli.Search.WithBody(bytes.NewReader(q)),
	)
	if err != nil {
		r.m.IncDocumentRepositoryRequests(op, esRes.Status(), time.Since(t).Seconds())
		return nil, fmt.Errorf("search documents: %w", err)
	}

	defer esRes.Body.Close()

	r.m.IncDocumentRepositoryRequests(op, esRes.Status(), time.Since(t).Seconds())

	if err := elasticutil.IsErr(esRes); err != nil {
		return nil, fmt.Errorf("search error: %s: %w", esRes.Status(), mapErrorToModel(err))
	}

	raw, err := elasticutil.ParseSearchResponse(esRes.Body)
	if err != nil {
		return nil, fmt.Errorf("parse search response: %w", err)
	}

	res, err := mapSearchResponseToDocuments(raw)
	if err != nil {
		return nil, fmt.Errorf("map search response to documents: %w", err)
	}

	return res, nil
}
