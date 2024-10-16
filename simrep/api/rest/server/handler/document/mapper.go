package document

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	openapi "simrep/api/rest/gen"
	"simrep/internal/model"
	"time"
)

var (
	errNilPart    = errors.New("nil part")
	errNoDocument = errors.New("no document")
	errNoBody     = errors.New("no body")
)

func mapUploadRequestToCommand(
	in openapi.PostDocumentUploadRequestObject,
) (model.DocumentFileUploadCommand, error) {
	var (
		file   model.File
		parsed bool
	)

	for {
		part, err := in.Body.NextPart()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return model.DocumentFileUploadCommand{}, fmt.Errorf("read part: %w", err)
		}

		if part.FormName() != "document" {
			continue
		}

		item, err := mapPartToFile(part)
		if err != nil {
			return model.DocumentFileUploadCommand{}, fmt.Errorf("map part to file: %w", err)
		}

		file = *item
		parsed = true

		break
	}

	if !parsed {
		return model.DocumentFileUploadCommand{}, errNoDocument
	}

	return model.DocumentFileUploadCommand{
		Item: file,
	}, nil
}

func mapPartToFile(in *multipart.Part) (*model.File, error) {
	if in == nil {
		return nil, errNilPart
	}

	buf := &bytes.Buffer{}

	_, err := io.Copy(buf, in)
	if err != nil {
		return nil, fmt.Errorf("read part: %w", err)
	}

	return &model.File{
		Name:        in.FileName(),
		Content:     buf.Bytes(),
		Sha256:      sha256String(buf.Bytes()),
		LastUpdated: time.Time{},
	}, nil
}

func mapUploadCommandToResponse(
	cmd model.DocumentFileUploadCommand,
) openapi.PostDocumentUpload200JSONResponse {
	return openapi.PostDocumentUpload200JSONResponse{
		UploadSuccessJSONResponse: openapi.UploadSuccessJSONResponse(
			openapi.UploadSuccess{
				DocumentID: &cmd.Item.Sha256,
			},
		),
	}
}

func mapDocumentFileRequestToQuery(
	in openapi.GetDocumentDocumentIdDownloadRequestObject,
) model.DocumentFileQuery {
	return model.DocumentFileQuery{
		ID: in.DocumentId,
	}
}

func mapSearchRequestToQuery(
	in openapi.PostDocumentsSearchRequestObject,
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

func mapDocumentsToSearchResponse(in []model.Document) openapi.PostDocumentsSearch200JSONResponse {
	docs := make([]openapi.DocumentSummary, 0, len(in))

	for _, v := range in {
		v := v

		docs = append(docs, openapi.DocumentSummary{
			LastUpdated: &v.LastUpdated,
			Id:          &v.ID,
			Name:        &v.Name,
		})
	}

	return openapi.PostDocumentsSearch200JSONResponse{
		SearchResultJSONResponse: openapi.SearchResultJSONResponse{
			Documents: &docs,
		},
	}
}

func mapFileToDownloadResponse(
	in model.DocumentFile,
) openapi.GetDocumentDocumentIdDownload200ApplicationoctetStreamResponse {
	return openapi.GetDocumentDocumentIdDownload200ApplicationoctetStreamResponse{
		Body: bytes.NewReader(in.Content),
		Headers: openapi.GetDocumentDocumentIdDownload200ResponseHeaders{
			ContentDisposition: fmt.Sprintf(`attachment; filename="%s"`, in.Name),
		},
		ContentLength: int64(len(in.Content)),
	}
}

func sha256String(in []byte) string {
	hash := sha256.New()
	_, _ = hash.Write(in)

	return hex.EncodeToString(hash.Sum(nil))
}
