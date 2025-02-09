package logging

import (
	"context"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go/micro"
	"github.com/zuzuka28/simreport/prj/fulltextindex/internal/model"
	pb "github.com/zuzuka28/simreport/prj/similarityindex/pkg/pb/v1"
)

func NewMiddleware() pb.Middleware {
	return func(h pb.Handler) pb.Handler {
		return func(ctx context.Context, req micro.Request) {
			lreq := newLoggedRequest(req)

			h(ctx, lreq)

			attrs := []any{
				"request_id", ctx.Value(model.RequestIDKey),
				"status", lreq.Status,
				"response_size_bytes", lreq.ResponseSize,
				"elapsed_time", time.Since(lreq.Start),
			}

			if lreq.ErrorDescription != "" {
				attrs = append(attrs, "error", lreq.ErrorDescription)
			}

			slog.Info(
				"request processed",
				attrs...,
			)
		}
	}
}

type loggedRequest struct {
	micro.Request
	Status           string
	ResponseSize     int
	Start            time.Time
	ErrorDescription string
}

func newLoggedRequest(msg micro.Request) *loggedRequest {
	return &loggedRequest{
		Request:          msg,
		Status:           "",
		ResponseSize:     0,
		Start:            time.Now(),
		ErrorDescription: "",
	}
}

func (m *loggedRequest) Respond(data []byte, opts ...micro.RespondOpt) error {
	m.Status = "200"
	m.ResponseSize = len(data)

	return m.Request.Respond(data, opts...) //nolint:wrapcheck
}

func (m *loggedRequest) Error(code, description string, data []byte, opts ...micro.RespondOpt) error {
	m.Status = code
	m.ResponseSize = len(data)
	m.ErrorDescription = string(data)

	return m.Request.Error(code, description, data, opts...) //nolint:wrapcheck
}
