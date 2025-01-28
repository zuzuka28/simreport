package documentstatus

import (
	"context"
	"document/internal/model"
	"fmt"
)

func (s *Service) Update(
	ctx context.Context,
	cmd model.DocumentStatusUpdateCommand,
) error {
	if err := s.r.Update(ctx, cmd); err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	return nil
}
