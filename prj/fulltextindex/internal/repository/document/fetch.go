package document

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/fulltextindex/internal/model"
)

func (s *Repository) Fetch(
	ctx context.Context,
	query model.DocumentQuery,
) (model.Document, error) {
	reqbody := []byte(query.ID)

	resp, err := s.conn.RequestWithContext(ctx, s.endpoint(s.endpointByID), reqbody)
	if err != nil {
		return model.Document{}, fmt.Errorf("do request: %w", err)
	}

	if err := isErr(resp); err != nil {
		return model.Document{}, err
	}

	res, err := parseFetchDocumentResponse(resp.Data)
	if err != nil {
		return model.Document{}, fmt.Errorf("parse fetch document: %w", err)
	}

	return res, nil
}

func (s *Repository) endpoint(
	method string,
) string {
	return s.group + "." + method
}
