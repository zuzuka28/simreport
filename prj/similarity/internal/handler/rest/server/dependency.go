package server

import (
	"context"

	openapi "github.com/zuzuka28/simreport/prj/similarity/internal/handler/rest/gen"
)

type (
	Metrics interface {
		IncHTTPRequest(op string, status string, size int, dur float64)
	}

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
