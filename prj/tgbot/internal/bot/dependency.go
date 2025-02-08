package bot

import (
	"context"

	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"
)

type (
	Metrics interface {
		IncBotRequestsByUser(username, userid string)
		IncBotErrors(desc string)
	}

	UserStateService interface {
		Fetch(
			ctx context.Context,
			query model.UserStateQuery,
		) (*model.UserState, error)
		Save(
			ctx context.Context,
			cmd model.UserStateSaveCommand,
		) error
	}

	DocumentService interface {
		Save(
			ctx context.Context,
			cmd model.DocumentSaveCommand,
		) (*model.DocumentSaveResult, error)
	}

	SimilarityService interface {
		Search(
			ctx context.Context,
			query model.SimilarityQuery,
		) ([]*model.SimilarityMatch, error)
	}
)
