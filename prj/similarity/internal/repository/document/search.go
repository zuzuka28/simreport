package document

import (
	"context"
	"fmt"
	"time"

	"github.com/zuzuka28/simreport/prj/similarity/internal/model"
)

func (s *Repository) Search(
	ctx context.Context,
	query model.DocumentSearchQuery,
) ([]model.Document, error) {
	const op = "search"

	t := time.Now()

	resp, err := s.cli.SearchDocument(ctx, mapDocumentSearchQueryToPb(query))
	if err != nil {
		s.m.IncDocumentRepositoryRequests(op, mapErrorToStatus(err), time.Since(t).Seconds())
		return nil, fmt.Errorf("do request: %w", mapErrorToModel(err))
	}

	s.m.IncDocumentRepositoryRequests(op, "200", time.Since(t).Seconds())

	res, err := parseSearchDocumentsResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("parse search documents: %w", err)
	}

	return res, nil
}
