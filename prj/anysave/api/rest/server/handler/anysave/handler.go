package anysave

import (
	openapi "anysave/api/rest/gen"
	"context"
	"fmt"
)

type Handler struct {
	s Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		s: s,
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
