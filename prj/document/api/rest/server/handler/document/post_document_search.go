package document

import (
	"context"
	"fmt"

	openapi "github.com/zuzuka28/simreport/prj/document/api/rest/gen"
)

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
