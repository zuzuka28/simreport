package analyze

import (
	"context"
	"simrep/internal/model"
)

type (
	VectorizerService interface {
		TextToVector(
			ctx context.Context,
			params model.VectorizeTextParams,
		) (model.Vector, error)
		ImageToVector(
			ctx context.Context,
			params model.VectorizeImageParams,
		) (model.Vector, error)
		ImageToHashes(
			ctx context.Context,
			params model.HashImageParams,
		) (model.HashImage, error)
	}

	Repository interface {
		Save(
			ctx context.Context,
			cmd model.AnalyzedDocumentSaveCommand,
		) error
		SearchSimilar(
			ctx context.Context,
			query model.AnalyzedDocumentSimilarQuery,
		) ([]model.AnalyzedDocumentMatch, error)
		Fetch(
			ctx context.Context,
			query model.AnalyzedDocumentQuery,
		) (model.AnalyzedDocument, error)
	}
)
