package userstate

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/zuzuka28/simreport/lib/elasticutil"
	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"
)

//nolint:gochecknoglobals
var now = time.Now

func mapUserStateSaveCommandToInternal(
	us model.UserStateSaveCommand,
) *userState {
	return &userState{
		UserID:      us.Item.UserID,
		State:       us.Item.State,
		LastUpdated: now(),
	}
}

func parseFetchUserStateResponse(
	in *elasticutil.Hit,
) (*model.UserState, error) {
	var raw userState

	if err := json.Unmarshal(in.Source, &raw); err != nil {
		return nil, fmt.Errorf("unmarshal source: %w", err)
	}

	return &model.UserState{
		UserID: raw.UserID,
		State:  raw.State,
	}, nil
}

func mapErrorToModel(err error) error {
	switch {
	case errors.Is(err, elasticutil.ErrInvalid):
		return errors.Join(err, model.ErrInvalid)

	case errors.Is(err, elasticutil.ErrNotFound):
		return errors.Join(err, model.ErrNotFound)

	default:
		return err
	}
}
