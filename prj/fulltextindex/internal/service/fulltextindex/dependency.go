package fulltextindex

import (
	"context"

	"github.com/zuzuka28/simreport/prj/fulltextindex/internal/model"
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
