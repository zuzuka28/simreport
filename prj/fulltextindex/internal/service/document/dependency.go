package document

import (
	"context"

	"github.com/zuzuka28/simreport/prj/fulltextindex/internal/model"
)

type (
	Repository interface {
		Fetch(
			ctx context.Context,
			query model.DocumentQuery,
		) (model.Document, error)
	}
)
