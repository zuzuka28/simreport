package server

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"

	"github.com/zuzuka28/simreport/prj/shingleindex/api/nats/micro/middleware/metrics"
	pb "github.com/zuzuka28/simreport/prj/similarityindex/pkg/pb/v1"
)

const requestTimeout = 60 * time.Second

type Server struct {
	s *pb.SimilarityIndexServer
}

func NewServer(
	conn *nats.Conn,
	h Handler,
	m Metrics,
) *Server {
	compose := struct {
		pb.UnsafeSimilarityIndexServer
		Handler
	}{
		UnsafeSimilarityIndexServer: nil,
		Handler:                     h,
	}

	return &Server{
		s: pb.NewSimilarityIndexServer(
			pb.SimilarityIndexServerConfig{
				Config: micro.Config{
					Name:         "similarity_shingle",
					Endpoint:     nil,
					Version:      "0.0.1",
					Description:  "",
					Metadata:     nil,
					QueueGroup:   "",
					StatsHandler: nil,
					DoneHandler:  nil,
					ErrorHandler: nil,
				},
				RequestTimeout:       requestTimeout,
				Middleware:           metrics.NewMiddleware(m),
				RequestErrorHandler:  nil,
				ResponseErrorHandler: nil,
			},
			conn,
			compose,
		),
	}
}

func (s *Server) Start(ctx context.Context) error {
	return s.s.Start(ctx) //nolint:wrapcheck
}

func (s *Server) Stop() error {
	return s.s.Stop() //nolint:wrapcheck
}
