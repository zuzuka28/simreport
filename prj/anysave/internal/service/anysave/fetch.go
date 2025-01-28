// anysave - saves everything that is sent to the anysave service.
package anysave

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/anysave/internal/model"
)

func (s *Service) Fetch(
	ctx context.Context,
	query model.FileQuery,
) (model.File, error) {
	if query.Bucket == "" {
		query.Bucket = bucketAnysave
	}

	res, err := s.r.Fetch(ctx, query)
	if err != nil {
		return model.File{}, fmt.Errorf("fetch document file: %w", err)
	}

	return res, nil
}
