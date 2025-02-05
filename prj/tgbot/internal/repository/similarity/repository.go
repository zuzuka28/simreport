package similarity

import (
	"github.com/nats-io/nats.go"

	pb "github.com/zuzuka28/simreport/prj/similarity/pkg/pb/v1"
)

type Repository struct {
	cli *pb.SimilarityServiceClient
}

func NewRepository(conn *nats.Conn) *Repository {
	return &Repository{
		cli: pb.NewSimilarityServiceClient(
			pb.SimilarityServiceClientConfig{
				MicroSubject: "similarity",
			},
			conn,
		),
	}
}
