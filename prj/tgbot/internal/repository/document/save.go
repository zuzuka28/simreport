package document

import (
	"context"
	"fmt"
	"time"

	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"
)

func (s *Repository) Save(
	ctx context.Context,
	cmd model.DocumentSaveCommand,
) (*model.DocumentSaveResult, error) {
	const op = "save"

	t := time.Now()

	resp, err := s.cli.UploadDocument(ctx, mapDocumentSaveCommandToPb(cmd))
	if err != nil {
		s.m.IncDocumentRepositoryRequests(op, mapErrorToStatus(err), time.Since(t).Seconds())
		return nil, fmt.Errorf("do request: %w", mapErrorToModel(err))
	}

	s.m.IncDocumentRepositoryRequests(op, "200", time.Since(t).Seconds())

	return &model.DocumentSaveResult{
		ID: resp.GetDocumentId(),
	}, nil
}
