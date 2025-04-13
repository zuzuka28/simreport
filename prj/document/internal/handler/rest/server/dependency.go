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
		PostDocumentSearch(
			ctx context.Context,
			request openapi.PostDocumentSearchRequestObject,
		) (openapi.PostDocumentSearchResponseObject, error)
		PostDocumentUpload(
			ctx context.Context,
			request openapi.PostDocumentUploadRequestObject,
		) (openapi.PostDocumentUploadResponseObject, error)
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
