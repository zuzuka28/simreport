package semanticindex

import (
	"context"
	"document/internal/model"
)

type (
	Repository interface {
		SearchSimilar(
			ctx context.Context,
			query model.DocumentSimilarQuery,
		) ([]*model.DocumentSimilarMatch, error)
	}
)
