package documentsaved

import (
	"context"
	"simrep/internal/model"
)

type (
	DocumentService interface {
		Fetch(
			ctx context.Context,
			query model.DocumentQuery,
		) (model.Document, error)
	}

	AnalyzeService interface {
		Analyze(
			ctx context.Context,
			item model.Document,
		) (model.AnalyzedDocument, error)
		Save(
			ctx context.Context,
			cmd model.AnalyzedDocumentSaveCommand,
		) error
	}
)
