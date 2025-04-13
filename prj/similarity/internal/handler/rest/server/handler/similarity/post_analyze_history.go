package similarity

import (
	"context"
	"fmt"

	openapi "github.com/zuzuka28/simreport/prj/similarity/internal/handler/rest/gen"
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
