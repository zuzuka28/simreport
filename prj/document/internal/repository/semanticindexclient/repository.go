package semanticindexclient

import (
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
