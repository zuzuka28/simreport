package anysave

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/zuzuka28/simreport/prj/anysave/internal/model"

	"github.com/minio/minio-go/v7"
)

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
		var cerr minio.ErrorResponse
		if errors.As(err, &cerr) && cerr.StatusCode == http.StatusNotFound {
			return model.File{}, model.ErrNotFound
		}

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
