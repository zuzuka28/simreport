package fulltextindexclient

import (
	"errors"
	"fmt"

	"github.com/zuzuka28/simreport/prj/document/internal/model"

	pb "github.com/zuzuka28/simreport/prj/fulltextindex/pkg/pb/v1"
)

var errInternal = errors.New("internal error")

func parseSearchSimilarResponse(
	in *pb.SearchSimilarResponse,
) []*model.SimilarityMatch {
	if in == nil || in.GetMatches() == nil {
		return nil
	}

	items := make([]*model.SimilarityMatch, 0, len(in.GetMatches()))
	for _, v := range in.GetMatches() {
		items = append(items, mapDocumentSimilarMatchToModel(v))
	}

	return items
}

func mapDocumentSimilarMatchToModel(in *pb.SimilarityDocumentMatch) *model.SimilarityMatch {
	if in == nil {
		return nil
	}

	return &model.SimilarityMatch{
		ID:            in.GetId(),
		Rate:          in.GetRate(),
		Highlights:    in.GetHighlights(),
		SimilarImages: in.GetSimilarImages(),
	}
}

func isErr(in *pb.Error) error {
	if in == nil {
		return nil
	}

	status := in.GetStatus()
	switch status {
	case 404: //nolint:gomnd,mnd
		return fmt.Errorf("%w: %s", model.ErrNotFound, in.GetDescription())

	case 400: //nolint:gomnd,mnd
		return fmt.Errorf("%w: %s", model.ErrInvalid, in.GetDescription())

	default:
		return fmt.Errorf("%w: %s", errInternal, in.GetDescription())
	}
}
