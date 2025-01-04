package anysave

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"simrep/internal/model"

	"github.com/minio/minio-go/v7"
	"golang.org/x/sync/errgroup"
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

func (r *Repository) SaveMany(ctx context.Context, cmd model.FileSaveManyCommand) error {
	g, gCtx := errgroup.WithContext(ctx)

	for _, item := range cmd.Items {
		g.Go(func() error {
			return r.Save(gCtx, model.FileSaveCommand{
				Bucket: cmd.Bucket,
				Item:   item,
			})
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("save file: %w", err)
	}

	return nil
}

func (r *Repository) Save(ctx context.Context, cmd model.FileSaveCommand) error {
	opts := minio.PutObjectOptions{ //nolint:exhaustruct
		UserMetadata: map[string]string{},
	}

	if cmd.Item.Name != "" {
		opts.UserMetadata[userMetadataNameKey] = cmd.Item.Name
	}

	_, err := r.cli.PutObject(
		ctx,
		cmd.Bucket,
		cmd.Item.Sha256,
		bytes.NewReader(cmd.Item.Content),
		int64(len(cmd.Item.Content)),
		minio.PutObjectOptions{ //nolint:exhaustruct
			UserMetadata: map[string]string{userMetadataNameKey: cmd.Item.Name},
		},
	)
	if err != nil {
		return fmt.Errorf("put file %s: %w", cmd.Item.Sha256, err)
	}

	return nil
}

func (r *Repository) Fetch(
	ctx context.Context,
	query model.FileQuery,
) (model.File, error) {
	objectInfo, err := r.cli.StatObject(
		ctx,
		query.Bucket,
		query.ID,
		minio.StatObjectOptions{}, //nolint:exhaustruct
	)
	if err != nil {
		return model.File{}, fmt.Errorf("stat object: %w", err)
	}

	object, err := r.cli.GetObject(
		ctx,
		query.Bucket,
		query.ID,
		minio.GetObjectOptions{}, //nolint:exhaustruct
	)
	if err != nil {
		return model.File{}, fmt.Errorf("get object: %w", err)
	}

	defer object.Close()

	var buf bytes.Buffer

	if _, err := io.Copy(&buf, object); err != nil {
		return model.File{}, fmt.Errorf("copy object data: %w", err)
	}

	return model.File{
		Name:        objectInfo.UserMetadata[userMetadataNameKey],
		Content:     buf.Bytes(),
		Sha256:      query.ID,
		LastUpdated: objectInfo.LastModified,
	}, nil
}
