package attribute

import (
	"context"
	openapi "document/api/rest/gen"
	"fmt"
)

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
