package document

import (
	"context"
	"simrep/internal/model"
)

type (
	Service interface {
		UploadFile(
			ctx context.Context,
			cmd model.DocumentFileUploadCommand,
		) error
		FetchFile(
			ctx context.Context,
			query model.DocumentFileQuery,
		) (model.DocumentFile, error)
		Search(
			ctx context.Context,
			query model.DocumentSearchQuery,
		) ([]model.Document, error)
	}
)
