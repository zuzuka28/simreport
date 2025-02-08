package document

import (
	"context"
	"fmt"
	"time"

	"github.com/zuzuka28/simreport/prj/similarity/internal/model"
)

func (s *Repository) Fetch(
	ctx context.Context,
	query model.DocumentQuery,
) (model.Document, error) {
	const op = "fetch"

	t := time.Now()

	resp, err := s.cli.FetchDocument(ctx, mapDocumentQueryToPb(query))
	if err != nil {
		s.m.IncDocumentRepositoryRequests(op, mapErrorToStatus(err), time.Since(t).Seconds())
		return model.Document{}, fmt.Errorf("do request: %w", mapErrorToModel(err))
	}

	s.m.IncDocumentRepositoryRequests(op, "200", time.Since(t).Seconds())

	res, err := parseFetchDocumentResponse(resp)
	if err != nil {
		return model.Document{}, fmt.Errorf("parse fetch document: %w", err)
	}

	return res, nil
}
