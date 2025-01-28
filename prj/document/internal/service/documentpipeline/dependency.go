package documentpipeline

import (
	"context"

	"github.com/zuzuka28/simreport/prj/document/internal/model"
)

type (
	Handler interface {
		Serve(ctx context.Context, documentID string) error
	}

	Stage struct {
		Trigger model.DocumentProcessingStatus
		Action  Handler
		Next    model.DocumentProcessingStatus
	}

	StatusService interface {
		Update(ctx context.Context, cmd model.DocumentStatusUpdateCommand) error
	}
)
