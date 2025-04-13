package metrics

import (
	"context"
	"time"

	"github.com/nats-io/nats.go/micro"
	pb "github.com/zuzuka28/simreport/prj/document/pkg/pb/v1"
)

func NewMiddleware(m Metrics) pb.Middleware {
	return func(h pb.Handler) pb.Handler {
		return func(ctx context.Context, req micro.Request) {
			h(ctx, newMeasuredRequest(req, m))
		}
	}
}

type measuredRequest struct {
	micro.Request
	m     Metrics
	start time.Time
}

func newMeasuredRequest(msg micro.Request, m Metrics) *measuredRequest {
	return &measuredRequest{
		Request: msg,
		m:       m,
		start:   time.Now(),
	}
}

func (m *measuredRequest) Respond(data []byte, opts ...micro.RespondOpt) error {
	m.m.IncNatsMicroRequest(m.Request.Subject(), "200", len(data), float64(time.Since(m.start).Seconds()))

	return m.Request.Respond(data, opts...) //nolint:wrapcheck
}

func (m *measuredRequest) Error(code, description string, data []byte, opts ...micro.RespondOpt) error {
	m.m.IncNatsMicroRequest(m.Request.Subject(), code, len(data), float64(time.Since(m.start).Seconds()))

	return m.Request.Error(code, description, data, opts...) //nolint:wrapcheck
}
