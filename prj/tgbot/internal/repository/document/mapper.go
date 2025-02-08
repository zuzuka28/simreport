package document

import (
	"errors"

	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"

	pb "github.com/zuzuka28/simreport/prj/document/pkg/pb/v1"
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

func mapDocumentSaveCommandToPb(
	in model.DocumentSaveCommand,
) *pb.UploadDocumentRequest {
	return &pb.UploadDocumentRequest{
		File: &pb.File{
			Content:  in.Item.Source.Content,
			Id:       in.Item.Source.Sha256,
			Filename: in.Item.Source.Name,
		},
		GroupIds: in.Item.GroupID,
		ParentId: in.Item.ParentID,
		Version:  int64(in.Item.Version),
	}
}
