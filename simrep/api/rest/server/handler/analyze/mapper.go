package analyze

import (
	openapi "simrep/api/rest/gen"
	"simrep/internal/model"
)

func mapSearchSimilarRequestToQuery(
	in openapi.GetAnalyzeDocumentIdSimilarRequestObject,
) model.DocumentSimilarQuery {
	return model.DocumentSimilarQuery{
		ID: in.DocumentId,
	}
}

func mapMatchesToSearchSimilarResponse(
	in []model.DocumentSimilarMatch,
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
