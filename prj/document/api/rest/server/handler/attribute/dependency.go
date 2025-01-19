package attribute

import (
	"context"
	"document/internal/model"
)

type (
	Service interface {
		Fetch(
			ctx context.Context,
			query model.AttributeQuery,
		) ([]model.Attribute, error)
	}
)
