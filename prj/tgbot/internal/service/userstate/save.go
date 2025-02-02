package userstate

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"
)

func (s *Service) Save(
	ctx context.Context,
	cmd model.UserStateSaveCommand,
) error {
	if err := s.r.Save(ctx, cmd); err != nil {
		return fmt.Errorf("save state: %w", err)
	}

	return nil
}
