package similarity

import (
	"context"
	openapi "simrep/api/rest/gen"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) PostFilesCompare(
	ctx context.Context,
	request openapi.PostFilesCompareRequestObject,
) (openapi.PostFilesCompareResponseObject, error) {
	return nil, nil
}

func (h *Handler) GetFilesFileIdCompareAll(
	ctx context.Context,
	request openapi.GetFilesFileIdCompareAllRequestObject,
) (openapi.GetFilesFileIdCompareAllResponseObject, error) {
	return nil, nil
}
