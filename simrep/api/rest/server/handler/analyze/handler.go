package analyze

import (
	"context"
	"fmt"
	openapi "simrep/api/rest/gen"
	"simrep/internal/model"
)

type Handler struct {
	s Service

	mapSearchSimilarRequest func(
		ctx context.Context,
		in openapi.PostAnalyzeSimilarRequestObject,
	) (model.AnalyzedDocumentSimilarQuery, error)
}

func NewHandler(s Service, dp DocumentParser) *Handler {
	return &Handler{
		s:                       s,
		mapSearchSimilarRequest: makeMapSearchSimilarRequestToQuery(s, dp),
	}
}

func (h *Handler) PostAnalyzeSimilar(
	ctx context.Context,
	params openapi.PostAnalyzeSimilarRequestObject,
) (openapi.PostAnalyzeSimilarResponseObject, error) {
	query, err := h.mapSearchSimilarRequest(ctx, params)
	if err != nil {
		return openapi.PostAnalyzeSimilar400JSONResponse{}, nil
	}

	res, err := h.s.SearchSimilar(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("search similar: %w", err)
	}

	return mapMatchesToSearchSimilarResponse(res), nil
}
