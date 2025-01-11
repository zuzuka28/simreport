package analyzehistory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"document/internal/model"
	"lib/elasticutil"

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

func (r *Repository) Save(
	ctx context.Context,
	cmd model.SimilarityHistorySaveCommand,
) error {
	cmd.Item.Date = now()

	raw := mapHistoryToInternal(cmd.Item)

	documentBytes, err := json.Marshal(raw)
	if err != nil {
		return fmt.Errorf("marshal doc: %w", err)
	}

	res, err := r.cli.Index(
		r.index,
		bytes.NewReader(documentBytes),
		r.cli.Index.WithDocumentID(cmd.Item.ID),
		r.cli.Index.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("index history: %w", err)
	}

	defer res.Body.Close()

	if err := elasticutil.IsErr(res); err != nil {
		return fmt.Errorf("save history: %w", mapErrorToModel(err))
	}

	return nil
}

func (r *Repository) Fetch(
	ctx context.Context,
	query model.SimilarityHistoryQuery,
) (*model.SimilarityHistoryList, error) {
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
		return nil, fmt.Errorf("search history: %w", err)
	}

	defer esRes.Body.Close()

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
