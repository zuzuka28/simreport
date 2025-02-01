package shingleindexclient

import (
	"github.com/nats-io/nats.go"

	pb "github.com/zuzuka28/simreport/prj/shingleindex/pkg/pb/v1"
)

type Repository struct {
	cli *pb.ShingleIndexServiceNatsClient
}

func NewRepository(conn *nats.Conn) *Repository {
	return &Repository{
		cli: pb.NewShingleIndexServiceClient(
			pb.ShingleIndexServiceNatsClientConfig{
				ServerName: "similarity_shingle",
			},
			conn,
		),
	}
}
