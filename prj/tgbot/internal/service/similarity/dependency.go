package similarity

import (
	"context"

	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"
)

type (
	Repository interface {
		Search(
			ctx context.Context,
			query model.SimilarityQuery,
		) ([]*model.SimilarityMatch, error)
	}
)
