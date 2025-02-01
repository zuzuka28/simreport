package document

import (
	"context"

	"github.com/zuzuka28/simreport/prj/document/internal/model"
)

type (
	Service interface {
		Search(
			ctx context.Context,
			query model.DocumentSearchQuery,
		) ([]model.Document, error)
		Save(
			ctx context.Context,
			cmd model.DocumentSaveCommand,
		) (*model.Document, error)
		Fetch(
			ctx context.Context,
			query model.DocumentQuery,
		) (model.Document, error)
	}

	StatusService interface {
		Update(
			ctx context.Context,
			cmd model.DocumentStatusUpdateCommand,
		) error
	}
)
