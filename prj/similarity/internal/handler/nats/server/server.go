package server

import (
	"context"
	"slices"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"

	"github.com/zuzuka28/simreport/prj/similarity/internal/handler/nats/middleware/logging"
	metricsmw "github.com/zuzuka28/simreport/prj/similarity/internal/handler/nats/middleware/metrics"
	"github.com/zuzuka28/simreport/prj/similarity/internal/handler/nats/middleware/reqid"
	pb "github.com/zuzuka28/simreport/prj/similarity/pkg/pb/v1"
)

const requestTimeout = 60 * time.Second

type Server struct {
	s *pb.SimilarityServiceServer
}

func NewServer(
	conn *nats.Conn,
	anh SimilarityHandler,
	m Metrics,
) *Server {
	compose := struct {
		pb.UnsafeSimilarityServiceServer
		SimilarityHandler
	}{
		UnsafeSimilarityServiceServer: nil,
		SimilarityHandler:             anh,
	}

	return &Server{
		s: pb.NewSimilarityServiceServer(
			pb.SimilarityServiceServerConfig{
				Config: micro.Config{
					Name:         "similarity",
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
				Middleware: useMiddleware(
					reqid.NewMiddleware(),
					metricsmw.NewMiddleware(m),
					logging.NewMiddleware(),
				),
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

func useMiddleware(mws ...pb.Middleware) pb.Middleware {
	return func(h pb.Handler) pb.Handler {
		for _, mw := range slices.Backward(mws) {
			h = mw(h)
		}

		return func(ctx context.Context, req micro.Request) {
			h(ctx, req)
		}
	}
}
