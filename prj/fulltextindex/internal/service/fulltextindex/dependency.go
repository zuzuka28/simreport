package fulltextindex

import (
	"context"
	"fulltextindex/internal/model"
)

type (
	Repository interface {
		Save(
			ctx context.Context,
			cmd model.DocumentSaveCommand,
		) error
		SearchSimilar(
			ctx context.Context,
			query model.DocumentSimilarQuery,
		) ([]model.DocumentSimilarMatch, error)
	}
)
