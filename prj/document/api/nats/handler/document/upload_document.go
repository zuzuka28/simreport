package document

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/document/internal/model"
	pb "github.com/zuzuka28/simreport/prj/document/pkg/pb/v1"
)

func (h *Handler) UploadDocument(
	ctx context.Context,
	params *pb.UploadRequest,
) (*pb.UploadDocumentResponse, error) {
	cmd := mapUploadRequestToCommand(params)

	doc, err := h.s.Save(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("upload document: %w", err)
	}

	if err := h.ss.Update(ctx, model.DocumentStatusUpdateCommand{
		ID:     doc.ID(),
		Status: model.DocumentProcessingStatusFileSaved,
	}); err != nil {
		return nil, fmt.Errorf("update document status: %w", err)
	}

	return mapUploadCommandToResponse(doc), nil
}
