package document

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/zuzuka28/simreport/prj/document/internal/model"

	"github.com/zuzuka28/simreport/lib/elasticutil"
)

func (r *Repository) Save(ctx context.Context, cmd model.DocumentSaveCommand) error {
	cmd.Item.LastUpdated = now()

	raw := mapDocumentToInternal(cmd.Item)

	documentBytes, err := json.Marshal(raw)
	if err != nil {
		return fmt.Errorf("marshal doc: %w", err)
	}

	res, err := r.cli.Index(
		r.index,
		bytes.NewReader(documentBytes),
		r.cli.Index.WithDocumentID(cmd.Item.ID()),
		r.cli.Index.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("index doc: %w", err)
	}

	defer res.Body.Close()

	if err := elasticutil.IsErr(res); err != nil {
		return fmt.Errorf("save document %s: %w", cmd.Item.ID(), mapErrorToModel(err))
	}

	return nil
}
