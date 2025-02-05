package semanticindexclient

import (
	"github.com/nats-io/nats.go"

	pb "github.com/zuzuka28/simreport/prj/similarityindex/pkg/pb/v1"
)

type Repository struct {
	cli *pb.SimilarityIndexClient
}

func NewRepository(conn *nats.Conn) *Repository {
	return &Repository{
		cli: pb.NewSimilarityIndexClient(
			pb.SimilarityIndexClientConfig{
				MicroSubject: "similarity_semantic",
			},
			conn,
		),
	}
}
