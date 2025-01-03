package anysave

import (
	"context"
	"fmt"
	openapi "simrep/api/rest/gen"
)

type Handler struct {
	s Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		s: s,
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
