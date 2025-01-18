package document

import (
	"context"
	"fmt"
	"shingleindex/internal/model"

	"github.com/nats-io/nats.go"
)

type Repository struct {
	conn         *nats.Conn
	group        string
	endpointByID string
}

func NewRepository(conn *nats.Conn) *Repository {
	return &Repository{
		conn:         conn,
		group:        "document",
		endpointByID: "byid",
	}
}

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
