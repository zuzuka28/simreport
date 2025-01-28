package analyze

import (
	"context"
	"fmt"

	openapi "github.com/zuzuka28/simreport/prj/document/api/rest/gen"
)

func (h *Handler) PostAnalyzeHistory(
	ctx context.Context,
	params openapi.PostAnalyzeHistoryRequestObject,
) (openapi.PostAnalyzeHistoryResponseObject, error) {
	query := mapAnalyzeHistoryRequestToQuery(params)

	res, err := h.s.SearchHistory(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("search history: %w", err)
	}

	return mapAnalyzeHistoryToResponse(res), nil
}
