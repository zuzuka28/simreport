package fulltextindexclient

import (
	"github.com/nats-io/nats.go"

	pb "github.com/zuzuka28/simreport/prj/fulltextindex/pkg/pb/v1"
)

type Repository struct {
	cli *pb.FullTextIndexServiceClient
}

func NewRepository(conn *nats.Conn) *Repository {
	return &Repository{
		cli: pb.NewFullTextIndexServiceClient(conn),
	}
}
