package minioutil

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewClientWithStartup(ctx context.Context, opts Config) (*minio.Client, error) {
	cli, err := minio.New(opts.Endpoint, &minio.Options{ //nolint:exhaustruct
		Creds: credentials.NewStaticV4(
			opts.AccessKeyID,
			opts.SecletAccessKey,
			"",
		),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("new s3 client: %w", err)
	}

	for _, bucket := range opts.Buckets {
		if err := setupBucket(ctx, cli, bucket); err != nil {
			return nil, fmt.Errorf("setup bucket: %w", err)
		}
	}

	return cli, nil
}

func setupBucket(
	ctx context.Context,
	cli *minio.Client,
	bucket string,
) error {
	exists, err := cli.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("check bucket exists: %w", err)
	}

	if exists {
		return nil
	}

	err = cli.MakeBucket(
		ctx,
		bucket,
		minio.MakeBucketOptions{}, //nolint:exhaustruct
	)
	if err != nil {
		return fmt.Errorf("create bucket %s: %w", bucket, err)
	}

	return nil
}
