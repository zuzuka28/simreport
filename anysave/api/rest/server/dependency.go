package server

import (
	openapi "anysave/api/rest/gen"
	"context"
)

type (
	FileHandler interface {
		PostUpload(
			ctx context.Context,
			request openapi.PostUploadRequestObject,
		) (openapi.PostUploadResponseObject, error)
		GetDocumentIdDownload(
			ctx context.Context,
			request openapi.GetDocumentIdDownloadRequestObject,
		) (openapi.GetDocumentIdDownloadResponseObject, error)
	}
)
