package anysave

import (
	"context"
	"simrep/internal/model"
)

type (
	Repository interface {
		Save(
			ctx context.Context,
			cmd model.FileSaveCommand,
		) error
		Fetch(
			ctx context.Context,
			query model.FileQuery,
		) (model.File, error)
	}
)
