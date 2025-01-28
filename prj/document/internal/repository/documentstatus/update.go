package documentstatus

import (
	"context"
	"document/internal/model"
	"fmt"
)

func (r *Repository) Update(
	ctx context.Context,
	cmd model.DocumentStatusUpdateCommand,
) error {
	key := cmd.ID

	fullsubject := subject + "." + string(cmd.Status)

	_, err := r.p.Publish(ctx, fullsubject, []byte(cmd.ID))
	if err != nil {
		return fmt.Errorf("publush update: %w", err)
	}

	_, err = r.kv.Put(ctx, key, []byte(cmd.Status))
	if err != nil {
		return fmt.Errorf("put status: %w", err)
	}

	return nil
}
