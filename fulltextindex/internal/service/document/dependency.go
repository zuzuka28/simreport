package document

import (
	"context"
	"fulltextindex/internal/model"
)

type (
	Repository interface {
		Fetch(
			ctx context.Context,
			query model.DocumentQuery,
		) (model.Document, error)
	}
)
