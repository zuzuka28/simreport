package document

import (
	"errors"
	openapi "simrep/api/rest/gen"
	"simrep/internal/model"
)

var errNoBody = errors.New("no body")

func mapSearchRequestToQuery(
	in openapi.PostDocumentSearchRequestObject,
) (model.DocumentSearchQuery, error) {
	if in.Body == nil {
		return model.DocumentSearchQuery{}, errNoBody
	}

	params := in.Body

	var name string
	if params.Name != nil {
		name = *params.Name
	}

	return model.DocumentSearchQuery{
		Name: name,
	}, nil
}

func mapDocumentsToSearchResponse(
	in []model.Document,
) openapi.PostDocumentSearch200JSONResponse {
	docs := make([]openapi.DocumentSummary, 0, len(in))

	for _, v := range in {
		v := v

		docs = append(docs, openapi.DocumentSummary{
			LastUpdated: &v.LastUpdated,
			Id:          &v.ID,
			Name:        &v.Name,
		})
	}

	return openapi.PostDocumentSearch200JSONResponse{
		SearchResultJSONResponse: openapi.SearchResultJSONResponse{
			Documents: &docs,
		},
	}
}
