package document

import (
	"context"
	"document/internal/model"
)

type (
	Service interface {
		Search(
			ctx context.Context,
			query model.DocumentSearchQuery,
		) ([]model.Document, error)
		Save(
			ctx context.Context,
			cmd model.DocumentSaveCommand,
		) (*model.Document, error)
	}
)
