package userstate

import (
	"context"

	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"
)

type Repository interface {
	Fetch(
		ctx context.Context,
		query model.UserStateQuery,
	) (*model.UserState, error)
	Save(
		ctx context.Context,
		cmd model.UserStateSaveCommand,
	) error
}
