package server

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
)

const (
	documentByIDSubject = "document.byid"
)

type Server struct {
	conn *nats.Conn
	dh   DocumentHandler
}

func NewServer(
	conn *nats.Conn,
	dh DocumentHandler,
) *Server {
	return &Server{
		conn: conn,
		dh:   dh,
	}
}

func (s *Server) Start(ctx context.Context) error {
	docbyid, err := s.conn.Subscribe(documentByIDSubject, s.dh.Fetch)
	if err != nil {
		return fmt.Errorf("subscribe document byid handler: %w", err)
	}

	defer func() { _ = docbyid.Drain() }()

	<-ctx.Done()

	return nil
}
