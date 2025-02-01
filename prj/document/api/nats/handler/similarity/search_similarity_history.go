package similarity

import (
	"context"
	"fmt"

	pb "github.com/zuzuka28/simreport/prj/document/pkg/pb/v1"
)

func (h *Handler) SearchSimilarityHistory(
	ctx context.Context,
	params *pb.SearchSimilarityHistoryRequest,
) (*pb.SearchSimilarityHistoryResponse, error) {
	q, err := mapSearchSimilarityHistoryRequestToModel(params)
	if err != nil {
		return nil, fmt.Errorf("map request to model: %w", err)
	}

	res, err := h.s.SearchHistory(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("search history: %w", err)
	}

	return mapSearchSimilarityHistoryResponseToPb(res), nil
}
