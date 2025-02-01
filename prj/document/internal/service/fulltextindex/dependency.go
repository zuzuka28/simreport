package fulltextindex

import (
	"context"

	"github.com/zuzuka28/simreport/prj/document/internal/model"
)

type (
	Repository interface {
		SearchSimilar(
			ctx context.Context,
			query model.SimilarityQuery,
		) ([]*model.SimilarityMatch, error)
	}
)
