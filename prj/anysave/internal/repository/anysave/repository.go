package anysave

import (
	"github.com/minio/minio-go/v7"
)

const userMetadataNameKey = "Name"

type Repository struct {
	cli *minio.Client
}

func NewRepository(
	s3 *minio.Client,
) *Repository {
	return &Repository{
		cli: s3,
	}
}
