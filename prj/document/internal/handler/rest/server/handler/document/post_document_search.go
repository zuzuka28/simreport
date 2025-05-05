package document

import (
	"context"
	"fmt"

	openapi "github.com/zuzuka28/simreport/prj/document/internal/handler/rest/gen"
)

func (h *Handler) PostSearch(
	ctx context.Context,
	params openapi.PostSearchRequestObject,
) (openapi.PostSearchResponseObject, error) {
	query, err := mapSearchRequestToQuery(params)
	if err != nil {
		return openapi.PostSearch400JSONResponse{}, nil
	}

	res, err := h.s.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("search documents: %w", err)
	}

	return mapDocumentsToSearchResponse(res), nil
}
