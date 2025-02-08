package similarityindexclient

import (
	"errors"

	"github.com/zuzuka28/simreport/prj/similarity/internal/model"

	pb "github.com/zuzuka28/simreport/prj/similarityindex/pkg/pb/v1"
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
