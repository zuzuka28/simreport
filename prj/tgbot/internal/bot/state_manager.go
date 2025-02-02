package bot

import (
	"context"
	"errors"
	"fmt"

	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"
)

type stateManager struct {
	uss UserStateService
}

func newStateManager(uss UserStateService) *stateManager {
	return &stateManager{uss: uss}
}

func (sm *stateManager) SwitchState(ctx context.Context, userID int, newState string) error {
	err := sm.uss.Save(ctx, model.UserStateSaveCommand{
		Item: model.UserState{
			UserID: userID,
			State:  newState,
		},
	})
	if err != nil {
		return fmt.Errorf("switch state: %w", err)
	}

	return nil
}

func (sm *stateManager) CurrentState(ctx context.Context, userID int) (string, error) {
	us, err := sm.uss.Fetch(ctx, model.UserStateQuery{UserID: userID})
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return string(botStateStart), nil
		}

		return "", fmt.Errorf("get current state: %w", err)
	}

	return us.State, nil
}
