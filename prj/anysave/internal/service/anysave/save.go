// anysave - saves everything that is sent to the anysave service.
package anysave

import (
	"anysave/internal/model"
	"context"
	"fmt"
)

func (s *Service) Save(
	ctx context.Context,
	cmd model.FileSaveCommand,
) error {
	if cmd.Bucket == "" {
		cmd.Bucket = bucketAnysave
	}

	_, err := s.r.Fetch(ctx, model.FileQuery{
		Bucket: cmd.Bucket,
		ID:     cmd.Item.Sha256,
	})
	if err == nil { // == nil
		return nil
	}

	if err := s.r.Save(ctx, cmd); err != nil {
		return fmt.Errorf("save file: %w", err)
	}

	if s.onSaveAction != nil {
		if err := s.onSaveAction(ctx, cmd); err != nil {
			return fmt.Errorf("on save action: %w", err)
		}
	}

	return nil
}
