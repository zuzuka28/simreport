package similarity

import (
	"errors"
	"fmt"

	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"

	pb "github.com/zuzuka28/simreport/prj/similarity/pkg/pb/v1"
)

var errInternal = errors.New("internal error")

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

func mapDocumentSearchQueryToPb(
	in model.SimilarityQuery,
) *pb.DocumentId {
	return &pb.DocumentId{
		DocumentId: in.ID,
	}
}

func mapSimilarityMatchesToModel(
	in *pb.SearchSimilarResponse,
) []*model.SimilarityMatch {
	res := make([]*model.SimilarityMatch, 0, len(in.GetDocuments()))

	for _, v := range in.GetDocuments() {
		res = append(res, &model.SimilarityMatch{
			ID:            v.GetId(),
			Rate:          float64(v.GetRate()),
			Highlights:    v.GetHighlights(),
			SimilarImages: v.GetSimilarImages(),
		})
	}

	return res
}
