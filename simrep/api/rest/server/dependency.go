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
		PostDocumentUpload(
			ctx context.Context,
			request openapi.PostDocumentUploadRequestObject,
		) (openapi.PostDocumentUploadResponseObject, error)
	}

	AnalyzeHandler interface {
		GetAnalyzeDocumentIdSimilar(
			ctx context.Context,
			request openapi.GetAnalyzeDocumentIdSimilarRequestObject,
		) (openapi.GetAnalyzeDocumentIdSimilarResponseObject, error)
		PostAnalyzeHistory(
			ctx context.Context,
			request openapi.PostAnalyzeHistoryRequestObject,
		) (openapi.PostAnalyzeHistoryResponseObject, error)
	}
)
