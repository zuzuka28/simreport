package document

import (
	"errors"
	"fmt"

	"github.com/zuzuka28/simreport/prj/shingleindex/internal/model"

	pb "github.com/zuzuka28/simreport/prj/document/pkg/pb/v1"
)

var errInternal = errors.New("internal error")

func parseFetchDocumentResponse(in *pb.FetchDocumentResponse) (model.Document, error) {
	raw := in.GetDocument()

	if in == nil || raw.GetText() == nil {
		return model.Document{}, nil
	}

	return model.Document{
		ID:   raw.GetText().GetId(),
		Text: raw.GetText().GetContent(),
	}, nil
}

func isErr(in *pb.Error) error {
	if in == nil {
		return nil
	}

	status := in.GetStatus()
	switch status {
	case 404:
		return fmt.Errorf("%w: %s", model.ErrNotFound, in.GetDescription())

	case 400:
		return fmt.Errorf("%w: %s", model.ErrInvalid, in.GetDescription())

	default:
		return fmt.Errorf("%w: %s", errInternal, in.GetDescription())
	}
}
