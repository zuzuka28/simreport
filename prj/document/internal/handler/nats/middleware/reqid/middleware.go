package reqid

import (
	"context"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go/micro"
	"github.com/zuzuka28/simreport/prj/document/internal/model"
	pb "github.com/zuzuka28/simreport/prj/document/pkg/pb/v1"
)

func NewMiddleware() pb.Middleware {
	return func(h pb.Handler) pb.Handler {
		return func(ctx context.Context, req micro.Request) {
			rid := req.Headers().Get(model.RequestIDHeader)

			if rid == "" {
				rid = uuid.NewString()
			}

			ctx = context.WithValue(ctx, model.RequestIDKey, rid)

			h(ctx, req)
		}
	}
}
