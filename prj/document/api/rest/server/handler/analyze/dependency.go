package analyze

import (
	"context"

	"github.com/zuzuka28/simreport/prj/document/internal/model"
)

type (
	Service interface {
		SearchSimilar(
			ctx context.Context,
			query model.DocumentSimilarQuery,
		) ([]*model.DocumentSimilarMatch, error)
		SearchHistory(
			ctx context.Context,
			query model.SimilarityHistoryQuery,
		) (*model.SimilarityHistoryList, error)
	}
)
