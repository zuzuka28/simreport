package document

import (
	"context"
	"fmt"

	pb "github.com/zuzuka28/simreport/prj/document/pkg/pb/v1"
)

func (h *Handler) SearchDocument(
	ctx context.Context,
	params *pb.SearchRequest,
) (*pb.SearchDocumentResponse, error) {
	q := mapSearchRequestToQuery(params)

	res, err := h.s.Search(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("search documents: %w", err)
	}

	return mapDocumentsToSearchResponse(res), nil
}
