package document

import (
	"context"
	"document/internal/model"
	"fmt"
)

func (s *Service) Parse(
	ctx context.Context,
	item model.File,
) (model.Document, error) {
	parsed, err := s.p.Parse(ctx, item)
	if err != nil {
		return model.Document{}, fmt.Errorf("parse document: %w", err)
	}

	return parsed, nil
}
