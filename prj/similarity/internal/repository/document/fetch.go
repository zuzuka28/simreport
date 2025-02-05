package document

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/similarity/internal/model"
)

func (s *Repository) Fetch(
	ctx context.Context,
	query model.DocumentQuery,
) (model.Document, error) {
	resp, err := s.cli.FetchDocument(ctx, mapDocumentQueryToPb(query))
	if err != nil {
		return model.Document{}, fmt.Errorf("do request: %w", err)
	}

	if err := isErr(resp.GetError()); err != nil {
		return model.Document{}, err
	}

	res, err := parseFetchDocumentResponse(resp)
	if err != nil {
		return model.Document{}, fmt.Errorf("parse fetch document: %w", err)
	}

	return res, nil
}
