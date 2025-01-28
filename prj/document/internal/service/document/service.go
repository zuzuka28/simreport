package document

import (
	"context"

	"github.com/zuzuka28/simreport/prj/document/internal/model"

	"github.com/google/uuid"
)

const (
	bucketText  = "texts"
	bucketImage = "images"
)

//nolint:gochecknoglobals
var genID = uuid.NewString

type Opts struct {
	OnSaveAction func(ctx context.Context, cmd model.DocumentSaveCommand) error
}

type Service struct {
	r  Repository
	fr FileRepository
	p  Parser

	onSaveAction func(ctx context.Context, cmd model.DocumentSaveCommand) error
}

func NewService(
	opts Opts,
	r Repository,
	fr FileRepository,
	p Parser,
) *Service {
	return &Service{
		r:            r,
		fr:           fr,
		p:            p,
		onSaveAction: opts.OnSaveAction,
	}
}
