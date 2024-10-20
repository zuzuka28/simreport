package document

import (
	"context"
	"simrep/internal/model"
)

type (
	Notify interface {
		Notify(
			ctx context.Context,
			documentID string,
			action model.NotifyAction,
			userdata any,
		) error
	}

	ImageRepository interface {
		Save(
			ctx context.Context,
			cmd model.FileSaveCommand,
		) error
		Fetch(
			ctx context.Context,
			query model.FileQuery,
		) (model.File, error)
	}

	FileRepository interface {
		Save(
			ctx context.Context,
			cmd model.FileSaveCommand,
		) error
		Fetch(
			ctx context.Context,
			query model.FileQuery,
		) (model.File, error)
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
