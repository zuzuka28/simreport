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

	DocumentService interface {
		Fetch(
			ctx context.Context,
			query model.DocumentQuery,
		) (model.Document, error)
	}

	ShingleIndexService interface {
		SearchSimilar(
			ctx context.Context,
			query model.DocumentSimilarQuery,
		) ([]model.DocumentSimilarMatch, error)
	}

	FulltextIndexService interface {
		SearchSimilar(
			ctx context.Context,
			query model.DocumentSimilarQuery,
		) ([]model.DocumentSimilarMatch, error)
	}

	Repository interface {
		Save(
			ctx context.Context,
			cmd model.AnalyzedDocumentSaveCommand,
		) error
		SearchSimilar(
			ctx context.Context,
			query model.DocumentSimilarQuery,
		) ([]model.DocumentSimilarMatch, error)
		Fetch(
			ctx context.Context,
			query model.AnalyzedDocumentQuery,
		) (model.AnalyzedDocument, error)
	}
)
