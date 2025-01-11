package document

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	openapi "document/api/rest/gen"
	"document/internal/model"
	"time"
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

func mapUploadRequestToCommand(
	in openapi.PostDocumentUploadRequestObject,
) (model.DocumentSaveCommand, error) {
	return model.DocumentSaveCommand{
		Item: model.Document{
			ID:          valOrDefault(in.Body.ParentID),
			Name:        "",
			LastUpdated: time.Time{},
			Version:     valOrDefault(in.Body.Version),
			GroupID:     []string{valOrDefault(in.Body.GroupID)},
			SourceID:    in.Body.FileID,
			TextID:      "",
			ImageIDs:    nil,
			WithContent: false,
			Source:      model.File{},
			Text:        model.File{},
			Images:      nil,
		},
	}, nil
}

func mapUploadCommandToResponse(
	doc *model.Document,
) openapi.PostDocumentUpload200JSONResponse {
	return openapi.PostDocumentUpload200JSONResponse{
		UploadSuccessJSONResponse: openapi.UploadSuccessJSONResponse(
			openapi.UploadSuccess{
				DocumentID: &doc.ID,
			},
		),
	}
}

func sha256String(in []byte) string {
	hash := sha256.New()
	_, _ = hash.Write(in)

	return hex.EncodeToString(hash.Sum(nil))
}

func valOrDefault[T any](in *T) T {
	if in == nil {
		var t T
		return t
	}

	return *in
}
