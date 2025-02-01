package server

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
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
	docsrv, err := micro.AddService(s.conn, micro.Config{ //nolint:exhaustruct
		Name:    "document",
		Version: "0.0.1",
	})
	if err != nil {
		return fmt.Errorf("create document service: %w", err)
	}

	docsrvg := docsrv.AddGroup("document")

	defer func() { _ = docsrv.Stop() }()

	if err := docsrvg.AddEndpoint("byid", micro.HandlerFunc(s.dh.Fetch)); err != nil {
		return fmt.Errorf("create document by id handler: %w", err)
	}

	<-ctx.Done()

	return nil
}
