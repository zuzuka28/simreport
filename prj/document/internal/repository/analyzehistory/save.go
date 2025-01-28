package analyzehistory

import (
	"bytes"
	"context"
	"document/internal/model"
	"encoding/json"
	"fmt"

	"github.com/zuzuka28/simreport/lib/elasticutil"
)

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
