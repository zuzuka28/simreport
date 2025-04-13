package similarity

import (
	"context"
	"fmt"

	pb "github.com/zuzuka28/simreport/prj/similarity/pkg/pb/v1"
)

func (h *Handler) SearchSimilar(
	ctx context.Context,
	params *pb.DocumentId,
) (*pb.SearchSimilarResponse, error) {
	q, err := mapSearchSimilarRequestToModel(params)
	if err != nil {
		return nil, fmt.Errorf("map request to model: %w", err)
	}

	res, err := h.s.SearchSimilar(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("search similar: %w", err)
	}

	return mapSearchSimilarResponseToPb(res), nil
}
