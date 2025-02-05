package fulltextindex

import (
	"context"

	"github.com/zuzuka28/simreport/prj/similarity/internal/model"
)

type (
	Repository interface {
		SearchSimilar(
			ctx context.Context,
			query model.SimilarityQuery,
		) ([]*model.SimilarityMatch, error)
	}
)
