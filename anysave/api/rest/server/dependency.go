package server

import (
	openapi "anysave/api/rest/gen"
	"context"
)

type (
	FileHandler interface {
		PostDocumentUpload(
			ctx context.Context,
			request openapi.PostDocumentUploadRequestObject,
		) (openapi.PostDocumentUploadResponseObject, error)
		GetDocumentDocumentIdDownload(
			ctx context.Context,
			request openapi.GetDocumentDocumentIdDownloadRequestObject,
		) (openapi.GetDocumentDocumentIdDownloadResponseObject, error)
	}
)
