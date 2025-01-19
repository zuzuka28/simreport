package anysave

import (
	"anysave/internal/model"
	"context"
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
