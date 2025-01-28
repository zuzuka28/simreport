package shingleindex

import (
	"context"

	"github.com/zuzuka28/simreport/prj/shingleindex/internal/model"
)

type (
	DocumentService interface {
		Fetch(
			ctx context.Context,
			query model.DocumentQuery,
		) (model.Document, error)
	}

	Repository interface {
		SearchSimilar(
			ctx context.Context,
			query model.MinhashSimilarQuery,
		) ([]*model.MinhashSimilarMatch, error)
		Save(
			ctx context.Context,
			cmd model.MinhashSaveCommand,
		) error
	}
)
