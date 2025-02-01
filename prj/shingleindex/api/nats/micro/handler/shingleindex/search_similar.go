package shingleindex

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/shingleindex/internal/model"
	pb "github.com/zuzuka28/simreport/prj/shingleindex/pkg/pb/v1"
)

func (h *Handler) SearchSimilar(
	ctx context.Context,
	params *pb.SearchSimilarRequest,
) (*pb.SearchSimilarResponse, error) {
	id := params.GetId()

	doc, err := h.ds.Fetch(ctx, model.DocumentQuery{
		ID: id,
	})
	if err != nil {
		return nil, fmt.Errorf("fetch source document: %w", err)
	}

	res, err := h.s.SearchSimilar(ctx, model.DocumentSimilarQuery{
		ID:   id,
		Item: doc,
	})
	if err != nil {
		return nil, fmt.Errorf("search similar: %w", err)
	}

	return mapDocumentToResponse(res), nil
}
