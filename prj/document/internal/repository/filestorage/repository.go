package filestorage

import (
	"github.com/minio/minio-go/v7"
)

const userMetadataNameKey = "Name"

const (
	bucketAnysave = "anysave"
)

type Repository struct {
	cli *minio.Client

	m Metrics
}

func NewRepository(
	s3 *minio.Client,
	m Metrics,
) *Repository {
	return &Repository{
		cli: s3,
		m:   m,
	}
}
