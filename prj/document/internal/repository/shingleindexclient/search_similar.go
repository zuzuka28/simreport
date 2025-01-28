package shingleindexclient

import (
	"context"
	"document/internal/model"
	"fmt"
)

func (s *Repository) SearchSimilar(
	ctx context.Context,
	query model.DocumentSimilarQuery,
) ([]*model.DocumentSimilarMatch, error) {
	reqbody := []byte(query.ID)

	resp, err := s.conn.RequestWithContext(ctx, s.endpoint(s.endpointSearch), reqbody)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	if err := isErr(resp); err != nil {
		return nil, err
	}

	res, err := parseSearchSimilarResponse(resp.Data)
	if err != nil {
		return nil, fmt.Errorf("parse search similar: %w", err)
	}

	return res, nil
}

func (s *Repository) endpoint(
	method string,
) string {
	return s.group + "." + method
}
