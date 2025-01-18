package shingleindex

import (
	"context"
	"shingleindex/internal/model"
)

type (
	Service interface {
		SearchSimilar(
			ctx context.Context,
			query model.DocumentSimilarQuery,
		) ([]*model.DocumentSimilarMatch, error)
	}

	DocumentService interface {
		Fetch(
			ctx context.Context,
			query model.DocumentQuery,
		) (model.Document, error)
	}
)
