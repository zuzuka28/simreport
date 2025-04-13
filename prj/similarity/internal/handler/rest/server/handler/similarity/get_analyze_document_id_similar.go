package similarity

import (
	"context"
	"fmt"

	openapi "github.com/zuzuka28/simreport/prj/similarity/internal/handler/rest/gen"
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
