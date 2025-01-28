package semanticindex

import (
	"context"

	"github.com/zuzuka28/simreport/prj/document/internal/model"
)

type (
	Repository interface {
		SearchSimilar(
			ctx context.Context,
			query model.DocumentSimilarQuery,
		) ([]*model.DocumentSimilarMatch, error)
	}
)
