package fulltextindex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/zuzuka28/simreport/lib/elasticutil"
	"github.com/zuzuka28/simreport/prj/fulltextindex/internal/model"
)

func (r *Repository) Save(
	ctx context.Context,
	cmd model.DocumentSaveCommand,
) error {
	const op = "save"

	raw := mapDocumentToInternal(cmd.Item)

	documentBytes, err := json.Marshal(raw)
	if err != nil {
		return fmt.Errorf("marshal doc: %w", err)
	}

	t := time.Now()

	esRes, err := r.cli.Index(
		r.index,
		bytes.NewReader(documentBytes),
		r.cli.Index.WithDocumentID(cmd.Item.ID),
		r.cli.Index.WithContext(ctx),
	)
	if err != nil {
		r.m.IncFulltextIndexRequests(op, metricsError, time.Since(t).Seconds())
		return fmt.Errorf("index doc: %w", err)
	}

	defer esRes.Body.Close()

	r.m.IncFulltextIndexRequests(op, strconv.Itoa(esRes.StatusCode), time.Since(t).Seconds())

	if err := elasticutil.IsErr(esRes); err != nil {
		return fmt.Errorf("save document %s: %w", cmd.Item.ID, err)
	}

	return nil
}
