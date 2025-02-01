package server

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"

	pb "github.com/zuzuka28/simreport/prj/fulltextindex/pkg/pb/v1"
)

const requestTimeout = 60 * time.Second

type Server struct {
	s *pb.FullTextIndexServiceNatsServer
}

func NewServer(
	conn *nats.Conn,
	h Handler,
) (*Server, error) {
	compose := struct {
		Handler
	}{
		Handler: h,
	}

	srv, err := pb.NewFullTextIndexServiceNatsServer(
		pb.FullTextIndexServiceNatsServerConfig{
			Config: micro.Config{
				Name:         "similarity_fulltext",
				Endpoint:     nil,
				Version:      "0.0.1",
				Description:  "",
				Metadata:     nil,
				QueueGroup:   "",
				StatsHandler: nil,
				DoneHandler:  nil,
				ErrorHandler: nil,
			},
			RequestTimeout: requestTimeout,
			Middleware:     nil,
			OnError:        nil,
		},
		conn,
		compose,
	)
	if err != nil {
		return nil, fmt.Errorf("create server: %w", err)
	}

	return &Server{
		s: srv,
	}, nil
}

func (s *Server) Stop() error {
	return s.s.Stop() //nolint:wrapcheck
}
