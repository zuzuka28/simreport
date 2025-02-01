package document

import (
	"context"
	"fmt"

	openapi "github.com/zuzuka28/simreport/prj/document/api/rest/gen"
)

//nolint:revive,stylecheck
func (h *Handler) GetDocumentIdDownload(
	ctx context.Context,
	params openapi.GetDocumentIdDownloadRequestObject,
) (openapi.GetDocumentIdDownloadResponseObject, error) {
	query := mapDocumentFileRequestToQuery(params)

	documentFile, err := h.s.Fetch(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("fetch file: %w", err)
	}

	return mapFileToDownloadResponse(documentFile.Source), nil
}
