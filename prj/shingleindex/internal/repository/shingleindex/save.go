package shingleindex

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/shingleindex/internal/model"
)

func (r *Repository) Save(
	ctx context.Context,
	cmd model.MinhashSaveCommand,
) error {
	if err := r.lsh.Insert(ctx, cmd.DocumentID, cmd.Minhash); err != nil {
		return fmt.Errorf("insert into lsh: %w", err)
	}

	return nil
}
