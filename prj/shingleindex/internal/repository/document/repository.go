package document

import (
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
