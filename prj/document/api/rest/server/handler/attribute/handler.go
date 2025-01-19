package attribute

import (
	"context"
	openapi "document/api/rest/gen"
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

func (h *Handler) PostAttribute(
	ctx context.Context,
	params openapi.PostAttributeRequestObject,
) (openapi.PostAttributeResponseObject, error) {
	query, err := mapAttributeRequestToQuery(params)
	if err != nil {
		return openapi.PostAttribute400JSONResponse{}, nil
	}

	res, err := h.s.Fetch(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("fetch attribute: %w", err)
	}

	return mapDocumentsToSearchResponse(res), nil
}
