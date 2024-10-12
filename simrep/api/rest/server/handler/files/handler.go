package files

import (
	"context"
	"fmt"
	openapi "simrep/api/rest/gen"
)

type Handler struct {
	s DocumentService
}

func NewHandler(s DocumentService) *Handler {
	return &Handler{
		s: s,
	}
}

func (h *Handler) PostFilesUpload(
	ctx context.Context,
	params openapi.PostFilesUploadRequestObject,
) (openapi.PostFilesUploadResponseObject, error) {
	cmd, err := mapUploadFilesToCommand(params)
	if err != nil {
		return nil, fmt.Errorf("map to upload command: %w", err)
	}

	if err := h.s.UploadManyFiles(ctx, cmd); err != nil {
		return nil, fmt.Errorf("upload files: %w", err)
	}

	return openapi.PostFilesUpload200JSONResponse{
		UploadSuccessJSONResponse: openapi.UploadSuccessJSONResponse(
			mapCommandToUploadSuccess(cmd),
		),
	}, nil
}

func (h *Handler) GetFilesFileId(
	ctx context.Context,
	request openapi.GetFilesFileIdRequestObject,
) (openapi.GetFilesFileIdResponseObject, error) {
	return nil, nil
}
