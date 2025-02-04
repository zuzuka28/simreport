package document

import (
	"context"
	"fmt"

	pb "github.com/zuzuka28/simreport/prj/document/pkg/pb/v1"
	"github.com/zuzuka28/simreport/prj/shingleindex/internal/model"
)

func (s *Repository) Fetch(
	ctx context.Context,
	query model.DocumentQuery,
) (model.Document, error) {
	resp, err := s.cli.FetchDocument(ctx, &pb.FetchDocumentRequest{
		Id:          query.ID,
		WithContent: true,
		Include: []pb.DocumentQueryInclude{
			pb.DocumentQueryInclude_DOCUMENT_QUERY_INCLUDE_TEXT,
		},
	})
	if err != nil {
		return model.Document{}, fmt.Errorf("do request: %w", err)
	}

	if err := isErr(resp.GetError()); err != nil {
		return model.Document{}, err
	}

	res, err := parseFetchDocumentResponse(resp)
	if err != nil {
		return model.Document{}, fmt.Errorf("parse fetch document: %w", err)
	}

	return res, nil
}
