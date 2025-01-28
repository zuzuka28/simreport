package anysave

import (
	"anysave/internal/model"
	"bytes"
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
)

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
