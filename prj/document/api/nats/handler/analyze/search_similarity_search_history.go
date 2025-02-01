package analyze

import (
	"context"
	"fmt"

	pb "github.com/zuzuka28/simreport/prj/document/pkg/pb/v1"
)

func (h *Handler) SearchSimilaritySearchHistory(
	ctx context.Context,
	params *pb.SearchSimilaritySearchHistoryRequest,
) (*pb.SearchSimilaritySearchHistoryResponse, error) {
	q, err := mapSearchSimilaritySearchHistoryRequestToModel(params)
	if err != nil {
		return nil, fmt.Errorf("map request to model: %w", err)
	}

	res, err := h.s.SearchHistory(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("search history: %w", err)
	}

	return mapSearchSimilaritySearchHistoryResponseToPb(res), nil
}
