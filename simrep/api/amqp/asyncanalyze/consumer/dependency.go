package consumer

import (
	"context"
	"simrep/internal/model"
)

type (
	DocumentService interface {
		FetchParsedFile(
			ctx context.Context,
			query model.ParsedDocumentFileQuery,
		) (model.ParsedDocumentFile, error)
	}

	AnalyzeService interface {
		Analyze(
			ctx context.Context,
			item model.ParsedDocumentFile,
		) (model.AnalyzedDocument, error)
		Save(
			ctx context.Context,
			cmd model.AnalyzedDocumentSaveCommand,
		) error
	}
)
