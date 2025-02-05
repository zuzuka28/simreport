package similarityindexclient

import (
	"github.com/nats-io/nats.go"

	pb "github.com/zuzuka28/simreport/prj/similarityindex/pkg/pb/v1"
)

type Opts struct {
	MicroSubject string
}

type Repository struct {
	cli *pb.SimilarityIndexClient
}

func NewRepository(
	cfg Opts,
	conn *nats.Conn,
) *Repository {
	return &Repository{
		cli: pb.NewSimilarityIndexClient(
			pb.SimilarityIndexClientConfig{
				MicroSubject: cfg.MicroSubject,
			},
			conn,
		),
	}
}
