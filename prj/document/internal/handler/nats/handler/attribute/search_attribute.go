package attribute

import (
	"context"
	"fmt"

	pb "github.com/zuzuka28/simreport/prj/document/pkg/pb/v1"
)

func (h *Handler) SearchAttribute(
	ctx context.Context,
	params *pb.SearchAttributeRequest,
) (*pb.SearchAttributeResponse, error) {
	q, err := mapSearchAttributeRequestToModel(params)
	if err != nil {
		return nil, fmt.Errorf("map request to model: %w", err)
	}

	res, err := h.s.Search(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("search attribute: %w", err)
	}

	return mapSearchAttributeToPb(res), nil
}
