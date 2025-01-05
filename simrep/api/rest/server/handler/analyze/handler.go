package analyze

import (
	"context"
	"fmt"
	openapi "simrep/api/rest/gen"
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
