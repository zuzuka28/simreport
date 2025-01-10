package documentstatus

import (
	"context"
	"anysave/internal/model"
)

type (
	Repository interface {
		Update(
			ctx context.Context,
			cmd model.DocumentStatusUpdateCommand,
		) error
		Fetch(
			ctx context.Context,
			query model.DocumentStatusQuery,
		) ([]*model.DocumentStatus, error)
	}
)
