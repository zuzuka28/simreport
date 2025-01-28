package analyze

import (
	"context"
	openapi "document/api/rest/gen"
	"fmt"
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
