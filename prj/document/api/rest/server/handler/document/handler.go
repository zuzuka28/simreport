package document

import (
	"context"
	"fmt"
	openapi "document/api/rest/gen"
)

type Handler struct {
	s Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		s: s,
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

	return mapUploadCommandToResponse(doc), nil
}
