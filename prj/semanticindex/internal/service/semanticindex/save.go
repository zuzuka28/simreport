package semanticindex

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/semanticindex/internal/model"
)

func (s *Service) Save(
	ctx context.Context,
	cmd model.DocumentSaveCommand,
) error {
	if cmd.Item.Vector == nil {
		vec, err := s.vs.TextToVector(ctx, model.VectorizeTextParams{
			Text: string(cmd.Item.Text),
		})
		if err != nil {
			return fmt.Errorf("vectorize document: %w", err)
		}

		cmd.Item.Vector = vec
	}

	if err := s.r.Save(ctx, cmd); err != nil {
		return fmt.Errorf("save analyzed document: %w", err)
	}

	return nil
}
