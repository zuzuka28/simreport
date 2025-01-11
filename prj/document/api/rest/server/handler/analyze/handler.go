package analyze

import (
	"context"
	"fmt"
	openapi "document/api/rest/gen"
)

type Handler struct {
	s Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		s: s,
	}
}

//nolint:revive
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
