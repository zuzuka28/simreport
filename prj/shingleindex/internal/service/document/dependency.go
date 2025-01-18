package document

import (
	"context"
	"shingleindex/internal/model"
)

type (
	Repository interface {
		Fetch(
			ctx context.Context,
			query model.DocumentQuery,
		) (model.Document, error)
	}
)
