package document

import (
	"context"
	"fmt"
	"shingleindex/internal/model"

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
		group:          "document.byid",
		endpointSearch: "search",
	}
}

func (s *Repository) Fetch(
	ctx context.Context,
	query model.DocumentQuery,
) (model.Document, error) {
	reqbody := []byte(query.ID)

	resp, err := s.conn.RequestWithContext(ctx, s.endpoint(s.endpointSearch), reqbody)
	if err != nil {
		return model.Document{}, fmt.Errorf("do request: %w", err)
	}

	if err := isErr(resp); err != nil {
		return model.Document{}, err
	}

	res, err := parseSearchSimilarResponse(resp.Data)
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
