package similarity

import (
	"context"

	"github.com/zuzuka28/simreport/prj/similarity/internal/model"
)

type (
	DocumentService interface {
		Fetch(
			ctx context.Context,
			query model.DocumentQuery,
		) (model.Document, error)
		Search(
			ctx context.Context,
			query model.DocumentSearchQuery,
		) ([]model.Document, error)
	}

	IndexingService interface {
		SearchSimilar(
			ctx context.Context,
			query model.SimilarityQuery,
		) ([]*model.SimilarityMatch, error)
	}

	Filestorage interface {
		Fetch(
			ctx context.Context,
			query model.FileQuery,
		) (model.File, error)
	}

	HistoryRepository interface {
		Save(
			ctx context.Context,
			cmd model.SimilarityHistorySaveCommand,
		) error
		Fetch(
			ctx context.Context,
			query model.SimilarityHistoryQuery,
		) (*model.SimilarityHistoryList, error)
	}
)
