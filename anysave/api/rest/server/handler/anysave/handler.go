package anysave

import (
	"context"
	"fmt"
	openapi "anysave/api/rest/gen"
	"anysave/internal/model"
)

type Handler struct {
	s  Service
	ss StatusService
}

func NewHandler(s Service, ss StatusService) *Handler {
	return &Handler{
		s:  s,
		ss: ss,
	}
}

func (h *Handler) PostDocumentUpload(
	ctx context.Context,
	params openapi.PostDocumentUploadRequestObject,
) (openapi.PostDocumentUploadResponseObject, error) {
	cmd, err := mapUploadRequestToCommand(params)
	if err != nil {
		return openapi.PostDocumentUpload400JSONResponse{}, nil
	}

	if err := h.s.Save(ctx, cmd); err != nil {
		return nil, fmt.Errorf("upload file: %w", err)
	}

	if err := h.ss.Update(ctx, model.DocumentStatusUpdateCommand{
		ID:     cmd.Item.Sha256,
		Status: model.DocumentProcessingStatusFileSaved,
	}); err != nil {
		return nil, fmt.Errorf("upload file: %w", err)
	}

	return mapUploadCommandToResponse(cmd), nil
}

func (h *Handler) GetDocumentDocumentIdDownload(
	ctx context.Context,
	params openapi.GetDocumentDocumentIdDownloadRequestObject,
) (openapi.GetDocumentDocumentIdDownloadResponseObject, error) {
	query := mapDocumentFileRequestToQuery(params)

	documentFile, err := h.s.Fetch(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("fetch file: %w", err)
	}

	return mapFileToDownloadResponse(documentFile), nil
}
