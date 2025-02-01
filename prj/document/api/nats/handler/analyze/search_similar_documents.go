package analyze

import (
	"context"
	"fmt"

	pb "github.com/zuzuka28/simreport/prj/document/pkg/pb/v1"
)

func (h *Handler) SearchSimilarDocuments(
	ctx context.Context,
	params *pb.DocumentId,
) (*pb.SearchSimilarDocumentsResponse, error) {
	q, err := mapSearchSimilarDocumentsRequestToModel(params)
	if err != nil {
		return nil, fmt.Errorf("map request to model: %w", err)
	}

	res, err := h.s.SearchSimilar(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("search similar: %w", err)
	}

	return mapSearchSimilarDocumentsResponseToPb(res), nil
}
