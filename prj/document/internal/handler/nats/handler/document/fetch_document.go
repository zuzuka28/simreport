package document

import (
	"context"
	"fmt"

	pb "github.com/zuzuka28/simreport/prj/document/pkg/pb/v1"
)

func (h *Handler) FetchDocument(
	ctx context.Context,
	params *pb.FetchDocumentRequest,
) (*pb.FetchDocumentResponse, error) {
	q := mapFetchDocumentRequestToQuery(params)

	res, err := h.s.Fetch(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("fetch document: %w", err)
	}

	return mapFetchDocumentResponseToPb(q.WithContent, res), nil
}
