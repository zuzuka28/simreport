package vectorizer

import (
	"context"

	"github.com/zuzuka28/simreport/prj/semanticindex/internal/model"
)

type (
	Repository interface {
		TextToVector(
			ctx context.Context,
			params model.VectorizeTextParams,
		) (model.Vector, error)
	}
)
