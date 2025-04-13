package attribute

import (
	"context"

	"github.com/zuzuka28/simreport/prj/document/internal/model"
)

type (
	Service interface {
		Search(
			ctx context.Context,
			query model.AttributeQuery,
		) ([]model.Attribute, error)
	}
)
