package document

import (
	"context"
	"simrep/internal/model"
)

type (
	Service interface {
		Search(
			ctx context.Context,
			query model.DocumentSearchQuery,
		) ([]model.Document, error)
	}
)
