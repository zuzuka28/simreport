package server

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"

	pb "github.com/zuzuka28/simreport/prj/document/pkg/pb/v1"
)

const requestTimeout = 60 * time.Second

type Server struct {
	s *pb.DocumentServiceNatsServer
}

func NewServer(
	conn *nats.Conn,
	doch DocumentHandler,
	attrh AttributeHandler,
	anh SimilarityHandler,
) (*Server, error) {
	compose := struct {
		DocumentHandler
		AttributeHandler
		SimilarityHandler
	}{
		DocumentHandler:   doch,
		AttributeHandler:  attrh,
		SimilarityHandler: anh,
	}

	srv, err := pb.NewDocumentServiceNatsServer(
		pb.DocumentServiceServerConfig{
			Config: micro.Config{
				Name:         "document",
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
