package server

import (
	"context"

	openapi "github.com/zuzuka28/simreport/prj/similarity/api/rest/gen"
)

type (
	SimilarityHandler interface {
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
