package document

import (
	"context"
	openapi "document/api/rest/gen"
	"document/internal/model"
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

func (h *Handler) PostDocumentSearch(
	ctx context.Context,
	params openapi.PostDocumentSearchRequestObject,
) (openapi.PostDocumentSearchResponseObject, error) {
	query, err := mapSearchRequestToQuery(params)
	if err != nil {
		return openapi.PostDocumentSearch400JSONResponse{}, nil
	}

	res, err := h.s.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("search documents: %w", err)
	}

	return mapDocumentsToSearchResponse(res), nil
}

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
