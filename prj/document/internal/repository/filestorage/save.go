package filestorage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/zuzuka28/simreport/prj/document/internal/model"

	"github.com/minio/minio-go/v7"
)

func (r *Repository) Save(ctx context.Context, cmd model.FileSaveCommand) error {
	if cmd.Bucket == "" {
		cmd.Bucket = bucketAnysave
	}

	_, err := r.cli.StatObject(
		ctx,
		cmd.Bucket,
		cmd.Item.Sha256,
		minio.StatObjectOptions{}, //nolint:exhaustruct
	)
	if err != nil {
		var cerr minio.ErrorResponse
		if errors.As(err, &cerr) && cerr.StatusCode != http.StatusNotFound {
			return fmt.Errorf("check object already exists: %w", err)
		}
	}

	opts := minio.PutObjectOptions{ //nolint:exhaustruct
		UserMetadata: map[string]string{},
	}

	if cmd.Item.Name != "" {
		opts.UserMetadata[userMetadataNameKey] = cmd.Item.Name
	}

	_, err = r.cli.PutObject(
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
