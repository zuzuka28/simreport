package document

import (
	"errors"
	"time"

	openapi "github.com/zuzuka28/simreport/prj/document/api/rest/gen"
	"github.com/zuzuka28/simreport/prj/document/internal/model"
)

var errNoBody = errors.New("no body")

func mapSearchRequestToQuery(
	in openapi.PostDocumentSearchRequestObject,
) (model.DocumentSearchQuery, error) {
	if in.Body == nil {
		return model.DocumentSearchQuery{}, errNoBody
	}

	return model.DocumentSearchQuery{
		GroupID:  valOrDefault(in.Body.GroupID),
		Name:     valOrDefault(in.Body.Name),
		ParentID: valOrDefault(in.Body.ParentID),
		Version:  valOrDefault(in.Body.Version),
	}, nil
}

func mapDocumentsToSearchResponse(
	in []model.Document,
) openapi.PostDocumentSearch200JSONResponse {
	docs := make([]openapi.DocumentSummary, 0, len(in))

	for _, v := range in {
		v := v

		docID := v.ID()

		docs = append(docs, openapi.DocumentSummary{
			GroupID:     &v.GroupID,
			Id:          &docID,
			LastUpdated: &v.LastUpdated,
			Name:        &v.Name,
			ParentID:    &v.ParentID,
			Version:     &v.Version,
		})
	}

	return openapi.PostDocumentSearch200JSONResponse{
		SearchResultJSONResponse: openapi.SearchResultJSONResponse{
			Documents: &docs,
		},
	}
}

func mapUploadRequestToCommand(
	in openapi.PostDocumentUploadRequestObject,
) (model.DocumentSaveCommand, error) { //nolint:unparam
	return model.DocumentSaveCommand{
		Item: model.Document{
			ParentID:    valOrDefault(in.Body.ParentID),
			Name:        "",
			LastUpdated: time.Time{},
			Version:     valOrDefault(in.Body.Version),
			GroupID:     valOrDefault(in.Body.GroupID),
			SourceID:    in.Body.FileID,
			TextID:      "",
			ImageIDs:    nil,
			WithContent: false,
			Source:      model.File{}, //nolint:exhaustruct
			Text:        model.File{}, //nolint:exhaustruct
			Images:      nil,
		},
	}, nil
}

func mapUploadCommandToResponse(
	doc *model.Document,
) openapi.PostDocumentUpload200JSONResponse {
	docID := doc.ID()

	return openapi.PostDocumentUpload200JSONResponse{
		UploadSuccessJSONResponse: openapi.UploadSuccessJSONResponse(
			openapi.UploadSuccess{
				DocumentID: &docID,
			},
		),
	}
}

func valOrDefault[T any](in *T) T {
	if in == nil {
		var t T
		return t
	}

	return *in
}
