package analyze

import (
	"context"
	openapi "document/api/rest/gen"
	"fmt"
)

//nolint:revive,stylecheck
func (h *Handler) GetAnalyzeDocumentIdSimilar(
	ctx context.Context,
	params openapi.GetAnalyzeDocumentIdSimilarRequestObject,
) (openapi.GetAnalyzeDocumentIdSimilarResponseObject, error) {
	query := mapSearchSimilarRequestToQuery(params)

	res, err := h.s.SearchSimilar(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("search similar: %w", err)
	}

	return mapMatchesToSearchSimilarResponse(res), nil
}
