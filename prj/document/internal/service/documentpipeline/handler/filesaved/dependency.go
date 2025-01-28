package filesaved

import (
	"context"

	"github.com/zuzuka28/simreport/prj/document/internal/model"
)

type (
	DocumentService interface {
		Fetch(
			ctx context.Context,
			query model.DocumentQuery,
		) (model.Document, error)
		Parse(
			_ context.Context,
			item model.File,
		) (model.Document, error)
		Save(
			ctx context.Context,
			cmd model.DocumentSaveCommand,
		) (*model.Document, error)
	}
)
