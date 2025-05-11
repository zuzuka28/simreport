package similarity

import (
	"time"

	"github.com/google/uuid"
)

//nolint:gochecknoglobals
var (
	now   = time.Now
	genID = uuid.NewString
)

type Opts struct{}

type Service struct {
	ds      DocumentService
	fs      Filestorage
	indices []IndexingService
	hr      HistoryRepository
}

func NewService(
	_ Opts,
	ds DocumentService,
	fs Filestorage,
	indices []IndexingService,
	hr HistoryRepository,
) *Service {
	return &Service{
		ds:      ds,
		fs:      fs,
		indices: indices,
		hr:      hr,
	}
}
