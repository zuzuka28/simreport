package similarityindexclient

import (
	"context"
	"log/slog"

	"github.com/nats-io/nats.go"

	"github.com/zuzuka28/simreport/prj/similarity/internal/model"
	pb "github.com/zuzuka28/simreport/prj/similarityindex/pkg/pb/v1"
)

type Opts struct {
	MicroSubject string
}

type Repository struct {
	index string
	cli   *pb.SimilarityIndexClient
	m     Metrics
}

func NewRepository(
	cfg Opts,
	conn *nats.Conn,
	m Metrics,
) *Repository {
	return &Repository{
		index: cfg.MicroSubject,
		cli: pb.NewSimilarityIndexClient(
			pb.SimilarityIndexClientConfig{
				MicroSubject: cfg.MicroSubject,
				Middleware:   middleware(cfg.MicroSubject),
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
