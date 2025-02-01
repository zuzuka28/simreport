package attribute

import (
	"context"

	"github.com/zuzuka28/simreport/prj/document/internal/model"
)

type (
	Repository interface {
		Search(
			ctx context.Context,
			query model.AttributeQuery,
		) ([]model.Attribute, error)
	}
)
