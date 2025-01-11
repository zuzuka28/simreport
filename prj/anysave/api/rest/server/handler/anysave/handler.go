package anysave

import (
	openapi "anysave/api/rest/gen"
	"anysave/internal/model"
	"context"
	"fmt"
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

func (h *Handler) PostUpload(
	ctx context.Context,
	params openapi.PostUploadRequestObject,
) (openapi.PostUploadResponseObject, error) {
	cmd, err := mapUploadRequestToCommand(params)
	if err != nil {
		return openapi.PostUpload400JSONResponse{}, nil
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

func (h *Handler) GetDocumentIdDownload(
	ctx context.Context,
	params openapi.GetDocumentIdDownloadRequestObject,
) (openapi.GetDocumentIdDownloadResponseObject, error) {
	query := mapDocumentFileRequestToQuery(params)

	documentFile, err := h.s.Fetch(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("fetch file: %w", err)
	}

	return mapFileToDownloadResponse(documentFile), nil
}
