package server

import (
	"context"
	openapi "simrep/api/rest/gen"
)

type (
	DocumentHandler interface {
		PostDocumentSearch(
			ctx context.Context,
			request openapi.PostDocumentSearchRequestObject,
		) (openapi.PostDocumentSearchResponseObject, error)
	}

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

	AnalyzeHandler interface {
		GetAnalyzeDocumentIdSimilar(
			ctx context.Context,
			request openapi.GetAnalyzeDocumentIdSimilarRequestObject,
		) (openapi.GetAnalyzeDocumentIdSimilarResponseObject, error)
	}
)
