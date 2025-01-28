package documentstatus

import (
	"context"

	"github.com/zuzuka28/simreport/prj/document/internal/model"
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
