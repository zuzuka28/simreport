package attribute

import (
	"context"
	"fmt"

	openapi "github.com/zuzuka28/simreport/prj/document/api/rest/gen"
)

func (h *Handler) PostAttribute(
	ctx context.Context,
	params openapi.PostAttributeRequestObject,
) (openapi.PostAttributeResponseObject, error) {
	query, err := mapAttributeRequestToQuery(params)
	if err != nil {
		return openapi.PostAttribute400JSONResponse{}, nil
	}

	res, err := h.s.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("search attribute: %w", err)
	}

	return mapDocumentsToSearchResponse(res), nil
}
