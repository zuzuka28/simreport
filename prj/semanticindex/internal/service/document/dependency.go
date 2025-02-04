package document

import (
	"context"

	"github.com/zuzuka28/simreport/prj/semanticindex/internal/model"
)

type (
	Repository interface {
		Fetch(
			ctx context.Context,
			query model.DocumentQuery,
		) (model.Document, error)
	}
)
