package anysave

import (
	"context"
	"anysave/internal/model"
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

	StatusService interface {
		Update(
			ctx context.Context,
			cmd model.DocumentStatusUpdateCommand,
		) error
	}
)
