package file

import (
	"context"
	"simrep/internal/model"
)

type (
	Service interface {
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
