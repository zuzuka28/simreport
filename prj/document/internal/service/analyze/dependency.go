package analyze

import (
	"context"

	"github.com/zuzuka28/simreport/prj/document/internal/model"
)

type (
	DocumentService interface {
		Fetch(
			ctx context.Context,
			query model.DocumentQuery,
		) (model.Document, error)
	}

	ShingleIndexService interface {
		SearchSimilar(
			ctx context.Context,
			query model.DocumentSimilarQuery,
		) ([]*model.DocumentSimilarMatch, error)
	}

	FulltextIndexService interface {
		SearchSimilar(
			ctx context.Context,
			query model.DocumentSimilarQuery,
		) ([]*model.DocumentSimilarMatch, error)
	}

	SemanticIndexService interface {
		SearchSimilar(
			ctx context.Context,
			query model.DocumentSimilarQuery,
		) ([]*model.DocumentSimilarMatch, error)
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
