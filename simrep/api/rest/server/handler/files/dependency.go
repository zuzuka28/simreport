package files

import (
	"context"
	"simrep/internal/model"
)

type (
	DocumentService interface {
		UploadManyFiles(
			ctx context.Context,
			cmd model.DocumentFileUploadManyCommand,
		) error
	}
)
