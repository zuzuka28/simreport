package semanticindex

import (
	"context"

	"github.com/zuzuka28/simreport/prj/semanticindex/internal/model"
)

type (
	VectorizerService interface {
		TextToVector(
			ctx context.Context,
			params model.VectorizeTextParams,
		) (model.Vector, error)
	}

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
