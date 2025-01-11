package filesaved

import (
	"context"
	"document/internal/model"
)

type (
	FileService interface {
		Fetch(
			ctx context.Context,
			query model.FileQuery,
		) (model.File, error)
	}

	DocumentService interface {
		Parse(
			_ context.Context,
			item model.File,
		) (model.Document, error)
		Save(
			ctx context.Context,
			cmd model.DocumentSaveCommand,
		) (*model.Document, error)
	}
)
