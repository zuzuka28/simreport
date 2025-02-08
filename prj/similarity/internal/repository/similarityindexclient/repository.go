package similarityindexclient

import (
	"github.com/nats-io/nats.go"

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
			},
			conn,
		),
		m: m,
	}
}
