package server

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
)

type Server struct {
	conn *nats.Conn
	dh   IndexerHandler
}

func NewServer(
	conn *nats.Conn,
	dh IndexerHandler,
) *Server {
	return &Server{
		conn: conn,
		dh:   dh,
	}
}

func (s *Server) Start(ctx context.Context) error {
	js, err := s.conn.JetStream()
	if err != nil {
		return fmt.Errorf("create jetstream: %w", err)
	}

	sub, err := js.QueueSubscribe("documentstatus.document_parsed", "fulltextindex_injest", s.dh.Save)
	if err != nil {
		return fmt.Errorf("create indexer sub: %w", err)
	}

	defer func() { _ = sub.Drain() }()

	<-ctx.Done()

	return nil
}
