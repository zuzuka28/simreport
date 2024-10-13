package image

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"simrep/internal/model"

	"github.com/minio/minio-go/v7"
	"golang.org/x/sync/errgroup"
)

type Opts struct {
	Bucket string `yaml:"bucket"`
}

type Repository struct {
	cli    *minio.Client
	bucket string
}

func NewRepository(
	opts Opts,
	s3 *minio.Client,
) *Repository {
	return &Repository{
		cli:    s3,
		bucket: opts.Bucket,
	}
}

func (r *Repository) SaveMany(ctx context.Context, cmd model.MediaFileSaveManyCommand) error {
	g, gCtx := errgroup.WithContext(ctx)

	for _, item := range cmd.Items {
		g.Go(func() error {
			return r.Save(gCtx, model.MediaFileSaveCommand{
				Item: item,
			})
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("save file: %w", err)
	}

	return nil
}

func (r *Repository) Save(ctx context.Context, cmd model.MediaFileSaveCommand) error {
	_, err := r.cli.PutObject(
		ctx,
		r.bucket,
		cmd.Item.Sha256,
		bytes.NewReader(cmd.Item.Content),
		int64(len(cmd.Item.Content)),
		minio.PutObjectOptions{}, //nolint:exhaustruct
	)
	if err != nil {
		return fmt.Errorf("put file %s: %w", cmd.Item.Sha256, err)
	}

	return nil
}

func (r *Repository) Fetch(
	ctx context.Context,
	query model.MediaFileQuery,
) (model.MediaFile, error) {
	objectInfo, err := r.cli.StatObject(
		ctx,
		r.bucket,
		query.ID,
		minio.StatObjectOptions{}, //nolint:exhaustruct
	)
	if err != nil {
		return model.MediaFile{}, fmt.Errorf("stat object: %w", err)
	}

	object, err := r.cli.GetObject(
		ctx,
		r.bucket,
		query.ID,
		minio.GetObjectOptions{}, //nolint:exhaustruct
	)
	if err != nil {
		return model.MediaFile{}, fmt.Errorf("get object: %w", err)
	}

	defer object.Close()

	var buf bytes.Buffer

	if _, err := io.Copy(&buf, object); err != nil {
		return model.MediaFile{}, fmt.Errorf("copy object data: %w", err)
	}

	return model.MediaFile{
		Content:     buf.Bytes(),
		Sha256:      query.ID,
		LastUpdated: objectInfo.LastModified,
	}, nil
}
