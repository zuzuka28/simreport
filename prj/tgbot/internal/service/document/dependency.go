package document

import (
	"context"

	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"
)

type (
	Repository interface {
		Save(
			ctx context.Context,
			cmd model.DocumentSaveCommand,
		) (*model.DocumentSaveResult, error)
	}
)
