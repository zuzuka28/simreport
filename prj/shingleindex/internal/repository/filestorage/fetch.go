package filestorage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/zuzuka28/simreport/prj/shingleindex/internal/model"

	"github.com/minio/minio-go/v7"
)

func (r *Repository) Fetch(
	ctx context.Context,
	query model.FileQuery,
) (model.File, error) {
	const op = "fetch"

	t := time.Now()

	if query.Bucket == "" {
		query.Bucket = bucketAnysave
	}

	objectInfo, err := r.cli.StatObject(
		ctx,
		query.Bucket,
		query.ID,
		minio.StatObjectOptions{}, //nolint:exhaustruct
	)
	if err != nil {
		var cerr minio.ErrorResponse
		if !errors.As(err, &cerr) {
			r.m.IncFilestorageRequests(op, metricsError, time.Since(t).Seconds())
			return model.File{}, fmt.Errorf("stat object: %w", err)
		}

		r.m.IncFilestorageRequests(op, strconv.Itoa(cerr.StatusCode), time.Since(t).Seconds())

		if cerr.StatusCode == http.StatusNotFound {
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
		r.m.IncFilestorageRequests(op, metricsError, time.Since(t).Seconds())
		return model.File{}, fmt.Errorf("get object: %w", err)
	}

	defer object.Close()

	var buf bytes.Buffer

	if _, err := io.Copy(&buf, object); err != nil {
		return model.File{}, fmt.Errorf("copy object data: %w", err)
	}

	r.m.IncFilestorageRequests(op, metricsSuccess, time.Since(t).Seconds())

	return model.File{
		Name:        objectInfo.UserMetadata[userMetadataNameKey],
		Content:     buf.Bytes(),
		Sha256:      query.ID,
		LastUpdated: objectInfo.LastModified,
	}, nil
}
