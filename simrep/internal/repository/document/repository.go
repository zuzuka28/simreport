package document

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"simrep/internal/model"
	"simrep/pkg/elasticutil"

	"github.com/elastic/go-elasticsearch/v8"
)

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

func (repo *Repository) SaveParsed(ctx context.Context, cmd model.ParsedDocumentSaveCommand) error {
	raw := mapParsedDocumentToInternal(cmd.Item)

	documentBytes, err := json.Marshal(raw)
	if err != nil {
		return fmt.Errorf("marshal doc: %w", err)
	}

	res, err := repo.cli.Index(
		repo.index,
		bytes.NewReader(documentBytes),
		repo.cli.Index.WithDocumentID(cmd.Item.ID),
		repo.cli.Index.WithContext(ctx),
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
