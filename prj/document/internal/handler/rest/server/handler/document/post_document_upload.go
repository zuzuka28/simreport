package document

import (
	"context"
	"fmt"

	openapi "github.com/zuzuka28/simreport/prj/document/internal/handler/rest/gen"
	"github.com/zuzuka28/simreport/prj/document/internal/model"
)

func (h *Handler) PostUpload(
	ctx context.Context,
	params openapi.PostUploadRequestObject,
) (openapi.PostUploadResponseObject, error) {
	cmd, err := mapUploadRequestToCommand(params)
	if err != nil {
		return openapi.PostUpload400JSONResponse{}, nil
	}

	doc, err := h.s.Save(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("upload document: %w", err)
	}

	if err := h.ss.Update(ctx, model.DocumentStatusUpdateCommand{
		ID:     doc.ID(),
		Status: model.DocumentProcessingStatusFileSaved,
	}); err != nil {
		return nil, fmt.Errorf("update document status: %w", err)
	}

	return mapUploadCommandToResponse(doc), nil
}
