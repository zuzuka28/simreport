package attribute

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
	query model.AttributeQuery,
) ([]model.Attribute, error) {
	const op = "search"

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
		r.m.IncAttributeRepositoryRequests(op, metricsError, time.Since(t).Seconds())
		return nil, fmt.Errorf("search attribute: %w", err)
	}

	defer esRes.Body.Close()

	r.m.IncAttributeRepositoryRequests(op, esRes.Status(), time.Since(t).Seconds())

	if err := elasticutil.IsErr(esRes); err != nil {
		return nil, fmt.Errorf("search error: %s: %w", esRes.Status(), mapErrorToModel(err))
	}

	raw, err := elasticutil.ParseSearchResponse(esRes.Body)
	if err != nil {
		return nil, fmt.Errorf("parse search response: %w", err)
	}

	res, err := mapSearchResponseToAttributes(raw)
	if err != nil {
		return nil, fmt.Errorf("map search response to attribute: %w", err)
	}

	return res, nil
}
