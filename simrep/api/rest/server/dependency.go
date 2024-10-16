package server

import (
	"context"
	openapi "simrep/api/rest/gen"
)

type (
	DocumentHandler interface {
		PostDocumentUpload(
			ctx context.Context,
			request openapi.PostDocumentUploadRequestObject,
		) (openapi.PostDocumentUploadResponseObject, error)
		GetDocumentDocumentIdDownload(
			ctx context.Context,
			request openapi.GetDocumentDocumentIdDownloadRequestObject,
		) (openapi.GetDocumentDocumentIdDownloadResponseObject, error)
		PostDocumentsSearch(
			ctx context.Context,
			request openapi.PostDocumentsSearchRequestObject,
		) (openapi.PostDocumentsSearchResponseObject, error)
	}

	AnalyzeHandler interface {
		PostAnalyzeSimilar(
			ctx context.Context,
			request openapi.PostAnalyzeSimilarRequestObject,
		) (openapi.PostAnalyzeSimilarResponseObject, error)
	}
)
