package analyze

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"simrep/internal/model"
	"simrep/pkg/elasticutil"

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
	cmd model.AnalyzedDocumentSaveCommand,
) error {
	raw := mapDocumentToInternal(cmd.Item)

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
		return fmt.Errorf("index doc: %w", err)
	}

	defer res.Body.Close()

	if err := elasticutil.IsErr(res); err != nil {
		return fmt.Errorf("save document %s: %w", cmd.Item.ID, err)
	}

	return nil
}

func (r *Repository) Fetch(
	ctx context.Context,
	query model.AnalyzedDocumentQuery,
) (model.AnalyzedDocument, error) {
	esRes, err := r.cli.Get(
		r.index,
		query.ID,
		r.cli.Get.WithContext(ctx),
	)
	if err != nil {
		return model.AnalyzedDocument{}, fmt.Errorf("fetch document %s: %w", query.ID, err)
	}

	defer esRes.Body.Close()

	if err := elasticutil.IsErr(esRes); err != nil {
		return model.AnalyzedDocument{}, fmt.Errorf("document %s not found: %w", query.ID, err)
	}

	raw, err := elasticutil.ParseDocResponse(esRes.Body)
	if err != nil {
		return model.AnalyzedDocument{}, fmt.Errorf("parse analyzed document response: %w", err)
	}

	res, err := parseDocument(raw)
	if err != nil {
		return model.AnalyzedDocument{}, fmt.Errorf("parse analyzed document: %w", err)
	}

	return res, nil
}

func (r *Repository) SearchSimilar(
	ctx context.Context,
	query model.AnalyzedDocumentSimilarQuery,
) ([]model.AnalyzedDocumentMatch, error) {
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

	res, err := mapSearchResponseToMatches(query, raw)
	if err != nil {
		return nil, fmt.Errorf("map search response to documents: %w", err)
	}

	return res, nil
}
