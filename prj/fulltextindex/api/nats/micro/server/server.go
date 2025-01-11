package server

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

type Server struct {
	conn *nats.Conn
	dh   FileindexHandler
}

func NewServer(
	conn *nats.Conn,
	dh FileindexHandler,
) *Server {
	return &Server{
		conn: conn,
		dh:   dh,
	}
}

func (s *Server) Start(ctx context.Context) error {
	docsrv, err := micro.AddService(s.conn, micro.Config{ //nolint:exhaustruct
		Name:    "similarity.fulltextindex",
		Version: "0.0.1",
	})
	if err != nil {
		return fmt.Errorf("create document service: %w", err)
	}

	docsrvg := docsrv.AddGroup("similarity.fulltextindex")

	defer func() { _ = docsrv.Stop() }()

	if err := docsrvg.AddEndpoint("search", micro.HandlerFunc(s.dh.SearchSimilar)); err != nil {
		return fmt.Errorf("create fulltext handler: %w", err)
	}

	<-ctx.Done()

	return nil
}
