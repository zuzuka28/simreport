package document

import (
	"context"
	"simrep/internal/model"
)

type (
	FileParser interface {
		Parse(
			ctx context.Context,
			item model.DocumentFile,
		) (model.ParsedDocumentFile, error)
	}

	ImageRepository interface {
		SaveMany(
			ctx context.Context,
			cmd model.MediaFileSaveManyCommand,
		) error
	}

	FileRepository interface {
		SaveMany(
			ctx context.Context,
			cmd model.MediaFileSaveManyCommand,
		) error
	}

	Repository interface {
		SaveParsed(
			ctx context.Context,
			cmd model.ParsedDocumentSaveCommand,
		) error
	}
)
