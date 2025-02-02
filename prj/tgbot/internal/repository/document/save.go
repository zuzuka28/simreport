package document

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"
)

func (s *Repository) Save(
	ctx context.Context,
	cmd model.DocumentSaveCommand,
) (*model.DocumentSaveResult, error) {
	resp, err := s.cli.UploadDocument(ctx, mapDocumentSaveCommandToPb(cmd))
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	if err := isErr(resp.GetError()); err != nil {
		return nil, err
	}

	return &model.DocumentSaveResult{
		ID: resp.GetDocumentId(),
	}, nil
}
