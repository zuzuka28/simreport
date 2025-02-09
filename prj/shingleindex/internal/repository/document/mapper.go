package document

import (
	"errors"

	"github.com/zuzuka28/simreport/prj/shingleindex/internal/model"

	pb "github.com/zuzuka28/simreport/prj/document/pkg/pb/v1"
)

var errInternal = errors.New("internal error")

func parseFetchDocumentResponse(in *pb.FetchDocumentResponse) model.Document {
	return model.Document{
		ID:       in.GetDocument().GetId(),
		SourceID: in.GetDocument().GetSource().GetId(),
		TextID:   in.GetDocument().GetText().GetId(),
		Text:     nil,
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

func mapDocumentQueryToPb(
	query model.DocumentQuery,
) *pb.FetchDocumentRequest {
	return &pb.FetchDocumentRequest{
		Id:          query.ID,
		WithContent: false,
		Include:     nil,
	}
}
