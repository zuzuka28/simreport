package document

import (
	"context"

	"github.com/zuzuka28/simreport/prj/document/internal/model"
)

type (
	Service interface {
		Fetch(
			ctx context.Context,
			query model.DocumentQuery,
		) (model.Document, error)
	}
)
