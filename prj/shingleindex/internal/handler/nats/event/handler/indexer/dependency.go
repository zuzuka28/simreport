package indexer

import (
	"context"

	"github.com/zuzuka28/simreport/prj/shingleindex/internal/model"
)

type (
	Service interface {
		Save(ctx context.Context, cmd model.DocumentSaveCommand) error
	}

	DocumentService interface {
		Fetch(
			ctx context.Context,
			query model.DocumentQuery,
		) (model.Document, error)
	}

	Filestorage interface {
		Fetch(
			ctx context.Context,
			query model.FileQuery,
		) (model.File, error)
	}
)
