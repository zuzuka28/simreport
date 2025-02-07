package document

import (
	"context"
	"fmt"
	"time"

	"github.com/zuzuka28/simreport/prj/document/internal/model"

	"github.com/zuzuka28/simreport/lib/elasticutil"
)

func (r *Repository) Fetch(ctx context.Context, query model.DocumentQuery) (model.Document, error) {
	const op = "fetch"

	t := time.Now()

	esRes, err := r.cli.Get(
		r.index,
		query.ID,
		r.cli.Get.WithContext(ctx),
	)
	if err != nil {
		r.m.IncDocumentRepositoryRequests(op, esRes.Status(), time.Since(t).Seconds())
		return model.Document{}, fmt.Errorf("fetch document %s: %w", query.ID, err)
	}

	defer esRes.Body.Close()

	r.m.IncDocumentRepositoryRequests(op, esRes.Status(), time.Since(t).Seconds())

	if err := elasticutil.IsErr(esRes); err != nil {
		return model.Document{}, fmt.Errorf("document %s not found: %w", query.ID, mapErrorToModel(err))
	}

	raw, err := elasticutil.ParseDocResponse(esRes.Body)
	if err != nil {
		return model.Document{}, fmt.Errorf("parse document response: %w", err)
	}

	res, err := parseDocument(raw)
	if err != nil {
		return model.Document{}, fmt.Errorf("parse document: %w", err)
	}

	return res, nil
}
