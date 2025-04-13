package similarity

import (
	openapi "github.com/zuzuka28/simreport/prj/similarity/internal/handler/rest/gen"
	"github.com/zuzuka28/simreport/prj/similarity/internal/model"
)

func mapSearchSimilarRequestToQuery(
	in openapi.GetAnalyzeDocumentIdSimilarRequestObject,
) model.SimilarityQuery {
	return model.SimilarityQuery{ //nolint:exhaustruct
		ID: in.DocumentId,
	}
}

func mapMatchesToSearchSimilarResponse(
	in []*model.SimilarityMatch,
) openapi.GetAnalyzeDocumentIdSimilarResponseObject {
	docs := make([]openapi.AnalyzedDocumentMatch, 0, len(in))

	for _, v := range in {
		v := v

		rate := float32(v.Rate)

		docs = append(docs, openapi.AnalyzedDocumentMatch{
			Highlights:    &v.Highlights,
			Id:            &v.ID,
			Rate:          &rate,
			SimilarImages: &v.SimilarImages,
		})
	}

	return openapi.GetAnalyzeDocumentIdSimilar200JSONResponse{
		SimilaritySearchResultJSONResponse: openapi.SimilaritySearchResultJSONResponse{
			Documents: &docs,
		},
	}
}

func mapAnalyzeHistoryRequestToQuery(
	params openapi.PostAnalyzeHistoryRequestObject,
) model.SimilarityHistoryQuery {
	body := params.Body

	return model.SimilarityHistoryQuery{
		DocumentID: valOrEmpty(body.DocumentID),
		Limit:      valOrEmpty(body.Limit),
		Offset:     valOrEmpty(body.Offset),
		DateFrom:   valOrEmpty(body.DateFrom),
		DateTo:     valOrEmpty(body.DateTo),
	}
}

func mapAnalyzeHistoryToResponse(
	in *model.SimilarityHistoryList,
) openapi.PostAnalyzeHistoryResponseObject {
	docs := make([]openapi.SimilaritySearchHistory, 0, len(in.Items))

	for _, v := range in.Items {
		matches := make([]openapi.AnalyzedDocumentMatch, 0, len(v.Matches))

		for _, m := range v.Matches {
			rate := float32(m.Rate)

			matches = append(matches, openapi.AnalyzedDocumentMatch{
				Highlights:    &m.Highlights,
				Id:            &m.ID,
				Rate:          &rate,
				SimilarImages: &m.SimilarImages,
			})
		}

		docs = append(docs, openapi.SimilaritySearchHistory{
			Date:       &v.Date,
			DocumentID: &v.DocumentID,
			Id:         &v.ID,
			Matches:    &matches,
		})
	}

	return openapi.PostAnalyzeHistory200JSONResponse{
		SimilaritySearchHistoryResultJSONResponse: openapi.SimilaritySearchHistoryResultJSONResponse{
			Count:     &in.Count,
			Documents: &docs,
		},
	}
}

func valOrEmpty[T any](v *T) T {
	if v == nil {
		var t T
		return t
	}

	return *v
}
