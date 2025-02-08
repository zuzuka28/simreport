package document

import (
	"context"
	"fmt"
	"time"

	"github.com/zuzuka28/simreport/prj/semanticindex/internal/model"
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

	return parseFetchDocumentResponse(resp), nil
}
