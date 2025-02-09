package document

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zuzuka28/simreport/prj/document/internal/model"

	"github.com/zuzuka28/simreport/lib/elasticutil"
)

func (r *Repository) Save(ctx context.Context, cmd model.DocumentSaveCommand) error {
	const op = "save"

	t := time.Now()

	cmd.Item.LastUpdated = now()

	raw := mapDocumentToInternal(cmd.Item)

	documentBytes, err := json.Marshal(raw)
	if err != nil {
		return fmt.Errorf("marshal doc: %w", err)
	}

	esRes, err := r.cli.Index(
		r.index,
		bytes.NewReader(documentBytes),
		r.cli.Index.WithDocumentID(cmd.Item.ID()),
		r.cli.Index.WithContext(ctx),
	)
	if err != nil {
		r.m.IncDocumentRepositoryRequests(op, metricsError, time.Since(t).Seconds())
		return fmt.Errorf("index doc: %w", err)
	}

	defer esRes.Body.Close()

	r.m.IncDocumentRepositoryRequests(op, esRes.Status(), time.Since(t).Seconds())

	if err := elasticutil.IsErr(esRes); err != nil {
		return fmt.Errorf("save document %s: %w", cmd.Item.ID(), mapErrorToModel(err))
	}

	return nil
}
