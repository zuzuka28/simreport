package analyzehistory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zuzuka28/simreport/prj/similarity/internal/model"

	"github.com/zuzuka28/simreport/lib/elasticutil"
)

func (r *Repository) Save(
	ctx context.Context,
	cmd model.SimilarityHistorySaveCommand,
) error {
	const op = "save"

	cmd.Item.Date = now()

	raw := mapHistoryToInternal(cmd.Item)

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
		r.m.IncAnalyzeHistoryRepositoryRequests(op, metricsError, time.Since(t).Seconds())
		return fmt.Errorf("index history: %w", err)
	}

	defer esRes.Body.Close()

	r.m.IncAnalyzeHistoryRepositoryRequests(op, esRes.Status(), time.Since(t).Seconds())

	if err := elasticutil.IsErr(esRes); err != nil {
		return fmt.Errorf("save history: %w", mapErrorToModel(err))
	}

	return nil
}
