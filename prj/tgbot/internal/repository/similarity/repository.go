package similarity

import (
	"context"
	"log/slog"

	"github.com/nats-io/nats.go"

	pb "github.com/zuzuka28/simreport/prj/similarity/pkg/pb/v1"
	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"
)

type Repository struct {
	cli *pb.SimilarityServiceClient
	m   Metrics
}

func NewRepository(
	conn *nats.Conn,
	m Metrics,
) *Repository {
	return &Repository{
		cli: pb.NewSimilarityServiceClient(
			pb.SimilarityServiceClientConfig{
				MicroSubject: "similarity",
				Middleware:   middleware("similarity"),
			},
			conn,
		),
		m: m,
	}
}

func middleware(sub string) func(pb.Invoker) pb.Invoker {
	return func(invoke pb.Invoker) pb.Invoker {
		return func(ctx context.Context, msg *nats.Msg) (*nats.Msg, error) {
			rid, ok := ctx.Value(model.RequestIDKey).(string)
			if ok {
				msg.Header.Set(model.RequestIDHeader, rid)
			}

			slog.Info(
				"send message",
				"subject", sub,
				"request_id", rid,
			)

			return invoke(ctx, msg)
		}
	}
}
