package userstate

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"
)

func (s *Service) Fetch(
	ctx context.Context,
	query model.UserStateQuery,
) (*model.UserState, error) {
	res, err := s.r.Fetch(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("fetch state: %w", err)
	}

	return res, nil
}
