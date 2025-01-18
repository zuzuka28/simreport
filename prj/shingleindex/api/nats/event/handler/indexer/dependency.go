package indexer

import (
	"context"
	"shingleindex/internal/model"
)

type (
	Service interface {
		Save(ctx context.Context, cmd model.DocumentSaveCommand) error
	}

	DocumentService interface {
		Fetch(
			ctx context.Context,
			query model.DocumentQuery,
		) (model.Document, error)
	}
)
