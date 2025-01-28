package fulltextindex

import (
	"context"
	"fmt"
	"fulltextindex/internal/model"
)

func (s *Service) Save(
	ctx context.Context,
	cmd model.DocumentSaveCommand,
) error {
	if err := s.r.Save(ctx, cmd); err != nil {
		return fmt.Errorf("save analyzed document: %w", err)
	}

	return nil
}
