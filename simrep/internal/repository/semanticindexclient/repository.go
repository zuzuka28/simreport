package fulltextindexclient

import (
	"context"
	"fmt"
	"simrep/internal/model"

	"github.com/nats-io/nats.go"
)

type Repository struct {
	conn           *nats.Conn
	group          string
	endpointSearch string
}

func NewRepository(conn *nats.Conn) *Repository {
	return &Repository{
		conn:           conn,
		group:          "similarity.semanticindex",
		endpointSearch: "search",
	}
}

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
