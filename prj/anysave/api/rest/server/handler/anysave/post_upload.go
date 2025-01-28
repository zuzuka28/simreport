package anysave

import (
	openapi "anysave/api/rest/gen"
	"context"
	"fmt"
)

func (h *Handler) PostUpload(
	ctx context.Context,
	params openapi.PostUploadRequestObject,
) (openapi.PostUploadResponseObject, error) {
	cmd, err := mapUploadRequestToCommand(params)
	if err != nil {
		return openapi.PostUpload400JSONResponse{}, nil
	}

	if err := h.s.Save(ctx, cmd); err != nil {
		return nil, fmt.Errorf("upload file: %w", err)
	}

	return mapUploadCommandToResponse(cmd), nil
}
