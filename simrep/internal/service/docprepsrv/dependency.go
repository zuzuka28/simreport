package docprepsrv

import (
	"context"
	"simrep/internal/model"
)

type (
	Repository interface {
		PreprocessRawDocument(
			ctx context.Context,
			doc []byte,
		) (*model.Document, error)
	}
)
