package fulltextindex

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/fulltextindex/internal/model"
	pb "github.com/zuzuka28/simreport/prj/similarityindex/pkg/pb/v1"
)

const bucketTexts = "texts"

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

	textfile, err := h.fs.Fetch(ctx, model.FileQuery{
		Bucket: bucketTexts,
		ID:     doc.TextID,
	})
	if err != nil {
		return nil, fmt.Errorf("fetch document text: %w", err)
	}

	doc.Text = textfile.Content

	res, err := h.s.SearchSimilar(ctx, model.DocumentSimilarQuery{
		ID:   id,
		Item: doc,
	})
	if err != nil {
		return nil, fmt.Errorf("search similar: %w", err)
	}

	return mapDocumentToResponse(res), nil
}
