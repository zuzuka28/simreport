package document

import (
	"context"
	"simrep/internal/model"
)

type (
	Service interface {
		Fetch(
			ctx context.Context,
			query model.DocumentQuery,
		) (model.Document, error)
	}
)
