package server

import (
	"context"

	openapi "github.com/zuzuka28/simreport/prj/anysave/api/rest/gen"
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
