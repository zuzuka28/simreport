package analyze

import (
	"context"
	"simrep/internal/model"
)

type (
	Service interface {
		Analyze(
			ctx context.Context,
			item model.Document,
		) (model.AnalyzedDocument, error)
		SearchSimilar(
			ctx context.Context,
			query model.AnalyzedDocumentSimilarQuery,
		) ([]model.AnalyzedDocumentMatch, error)
	}

	DocumentParser interface {
		Parse(
			_ context.Context,
			item model.File,
		) (model.Document, error)
	}
)
