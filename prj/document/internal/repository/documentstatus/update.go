package documentstatus

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/document/internal/model"
)

func (r *Repository) Update(
	ctx context.Context,
	cmd model.DocumentStatusUpdateCommand,
) error {
	key := cmd.ID

	fullsubject := subject + "." + string(cmd.Status)

	_, err := r.p.Publish(ctx, fullsubject, []byte(cmd.ID))
	if err != nil {
		r.m.IncDocumentStatusRepositoryUpdates(string(cmd.Status), metricsFailed)
		return fmt.Errorf("publush update: %w", err)
	}

	_, err = r.kv.Put(ctx, key, []byte(cmd.Status))
	if err != nil {
		r.m.IncDocumentStatusRepositoryUpdates(string(cmd.Status), metricsFailed)
		return fmt.Errorf("put status: %w", err)
	}

	r.m.IncDocumentStatusRepositoryUpdates(string(cmd.Status), metricsSuccess)

	return nil
}
