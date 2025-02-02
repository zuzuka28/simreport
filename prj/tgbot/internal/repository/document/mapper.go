package document

import (
	"errors"
	"fmt"

	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"

	pb "github.com/zuzuka28/simreport/prj/document/pkg/pb/v1"
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
