package document

import (
	"github.com/nats-io/nats.go"

	pb "github.com/zuzuka28/simreport/prj/document/pkg/pb/v1"
)

type Repository struct {
	cli *pb.DocumentServiceNatsClient
}

func NewRepository(conn *nats.Conn) *Repository {
	return &Repository{
		cli: pb.NewDocumentServiceClient(
			pb.DocumentServiceNatsClientConfig{
				ServerName: "document",
			},
			conn,
		),
	}
}
