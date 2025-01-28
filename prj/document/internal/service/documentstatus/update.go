package documentstatus

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/document/internal/model"
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
