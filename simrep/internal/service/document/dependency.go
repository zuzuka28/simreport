package document

import (
	"context"
	"simrep/internal/model"
)

type (
	ImageRepository interface {
		SaveMany(
			ctx context.Context,
			cmd model.MediaFileSaveManyCommand,
		) error
		Fetch(
			ctx context.Context,
			query model.MediaFileQuery,
		) (model.MediaFile, error)
	}

	FileRepository interface {
		SaveMany(
			ctx context.Context,
			cmd model.MediaFileSaveManyCommand,
		) error
		Fetch(
			ctx context.Context,
			query model.DocumentFileQuery,
		) (model.DocumentFile, error)
	}

	Repository interface {
		Save(
			ctx context.Context,
			cmd model.DocumentSaveCommand,
		) error
		Fetch(
			ctx context.Context,
			query model.DocumentQuery,
		) (model.Document, error)
		Search(
			ctx context.Context,
			query model.DocumentSearchQuery,
		) ([]model.Document, error)
	}
)
