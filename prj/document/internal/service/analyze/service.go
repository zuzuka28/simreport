package analyze

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
	ds         DocumentService
	shingleis  ShingleIndexService
	fulltextis FulltextIndexService
	semanticis SemanticIndexService
	hr         HistoryRepository
}

func NewService(
	_ Opts,
	ds DocumentService,
	shingleis ShingleIndexService,
	fulltextis FulltextIndexService,
	semanticis SemanticIndexService,
	hr HistoryRepository,
) *Service {
	return &Service{
		ds:         ds,
		shingleis:  shingleis,
		fulltextis: fulltextis,
		semanticis: semanticis,
		hr:         hr,
	}
}
