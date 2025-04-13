package similarity

import (
	"context"

	"github.com/zuzuka28/simreport/prj/similarity/internal/model"
)

type (
	Service interface {
		SearchSimilar(
			ctx context.Context,
			query model.SimilarityQuery,
		) ([]*model.SimilarityMatch, error)
		SearchHistory(
			ctx context.Context,
			query model.SimilarityHistoryQuery,
		) (*model.SimilarityHistoryList, error)
	}
)
