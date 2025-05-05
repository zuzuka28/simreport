package server

import (
	"context"

	openapi "github.com/zuzuka28/simreport/prj/document/internal/handler/rest/gen"
)

type (
	Metrics interface {
		IncHTTPRequest(op string, status string, size int, dur float64)
	}

	DocumentHandler interface {
		PostSearch(
			ctx context.Context,
			request openapi.PostSearchRequestObject,
		) (openapi.PostSearchResponseObject, error)
		PostUpload(
			ctx context.Context,
			request openapi.PostUploadRequestObject,
		) (openapi.PostUploadResponseObject, error)
		GetDocumentIdDownload(
			ctx context.Context,
			params openapi.GetDocumentIdDownloadRequestObject,
		) (openapi.GetDocumentIdDownloadResponseObject, error)
	}

	AttributeHandler interface {
		PostAttribute(
			ctx context.Context,
			request openapi.PostAttributeRequestObject,
		) (openapi.PostAttributeResponseObject, error)
	}
)
