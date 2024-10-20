package documentfile

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
