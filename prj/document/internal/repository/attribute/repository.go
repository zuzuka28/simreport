package attribute

import (
	"bytes"
	"context"
	"document/internal/model"
	"fmt"

	"github.com/zuzuka28/simreport/lib/elasticutil"

	"github.com/elastic/go-elasticsearch/v8"
)

//nolint:revive
type map_ map[string]any

type Opts struct {
	Index string `yaml:"index"`
}

type Repository struct {
	cli   *elasticsearch.Client
	index string
}

func NewRepository(
	opts Opts,
	es *elasticsearch.Client,
) (*Repository, error) {
	return &Repository{
		cli:   es,
		index: opts.Index,
	}, nil
}

func (r *Repository) Fetch(
	ctx context.Context,
	query model.AttributeQuery,
) ([]model.Attribute, error) {
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
		return nil, fmt.Errorf("search attribute: %w", err)
	}

	defer esRes.Body.Close()

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
