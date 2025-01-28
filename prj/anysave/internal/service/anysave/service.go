// anysave - saves everything that is sent to the anysave service.
package anysave

import (
	"context"

	"github.com/zuzuka28/simreport/prj/anysave/internal/model"
)

const (
	bucketAnysave = "github.com/zuzuka28/simreport/prj/anysave"
)

type Opts struct {
	OnSaveAction func(ctx context.Context, cmd model.FileSaveCommand) error
}

type Service struct {
	r            Repository
	onSaveAction func(ctx context.Context, cmd model.FileSaveCommand) error
}

func NewService(
	opts Opts,
	r Repository,
) *Service {
	return &Service{
		r:            r,
		onSaveAction: opts.OnSaveAction,
	}
}
