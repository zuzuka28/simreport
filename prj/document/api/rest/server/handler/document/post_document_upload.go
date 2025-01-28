package document

import (
	"context"
	openapi "document/api/rest/gen"
	"document/internal/model"
	"fmt"
)

func (h *Handler) PostDocumentUpload(
	ctx context.Context,
	params openapi.PostDocumentUploadRequestObject,
) (openapi.PostDocumentUploadResponseObject, error) {
	cmd, err := mapUploadRequestToCommand(params)
	if err != nil {
		return openapi.PostDocumentUpload400JSONResponse{}, nil
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
