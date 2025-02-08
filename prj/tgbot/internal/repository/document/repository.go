package document

import (
	"github.com/nats-io/nats.go"

	pb "github.com/zuzuka28/simreport/prj/document/pkg/pb/v1"
)

type Repository struct {
	cli *pb.DocumentServiceClient
	m   Metrics
}

func NewRepository(
	conn *nats.Conn,
	m Metrics,
) *Repository {
	return &Repository{
		cli: pb.NewDocumentServiceClient(
			pb.DocumentServiceClientConfig{
				MicroSubject: "document",
				Middleware:   nil,
			},
			conn,
		),
		m: m,
	}
}
