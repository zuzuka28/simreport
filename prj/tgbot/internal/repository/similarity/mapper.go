package similarity

import (
	"errors"

	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"

	pb "github.com/zuzuka28/simreport/prj/similarity/pkg/pb/v1"
)

var errInternal = errors.New("internal error")

func mapErrorToModel(err error) error {
	clierr := new(pb.ClientError)

	if !errors.As(err, &clierr) {
		return errors.Join(errInternal, err)
	}

	var merr error

	switch clierr.Status {
	case "404":
		merr = model.ErrNotFound

	case "400":
		merr = model.ErrInvalid

	default:
		merr = errInternal
	}

	return errors.Join(merr, err)
}

func mapErrorToStatus(err error) string {
	clierr := new(pb.ClientError)

	if !errors.As(err, &clierr) {
		return "500"
	}

	return clierr.Status
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
