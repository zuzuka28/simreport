package document

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	openapi "github.com/zuzuka28/simreport/prj/document/api/rest/gen"
	"github.com/zuzuka28/simreport/prj/document/internal/model"
)

var errNoBody = errors.New("no body")

var errNilPart = errors.New("nil part")

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

type documentUploadForm struct {
	item      model.Document
	callbacks map[string]func(part *multipart.Part) error
}

func newDocumentUploadForm() *documentUploadForm {
	p := &documentUploadForm{
		item:      model.Document{}, //nolint:exhaustruct
		callbacks: make(map[string]func(part *multipart.Part) error),
	}

	p.callbacks["parentID"] = p.parseParentID
	p.callbacks["version"] = p.parseVersion
	p.callbacks["groupID"] = p.parseGroupID
	p.callbacks["document"] = p.parseDocument

	return p
}

func (f *documentUploadForm) parseParentID(part *multipart.Part) error {
	data, err := io.ReadAll(part)
	if err != nil {
		return fmt.Errorf("reading parent_id: %w", err)
	}

	f.item.ParentID = string(data)

	return nil
}

func (f *documentUploadForm) parseVersion(part *multipart.Part) error {
	var buf bytes.Buffer

	if _, err := io.Copy(&buf, part); err != nil {
		return fmt.Errorf("reading version: %w", err)
	}

	ver, err := strconv.Atoi(buf.String())
	if err != nil {
		return fmt.Errorf("parsing version: %w", err)
	}

	f.item.Version = ver

	return nil
}

func (f *documentUploadForm) parseGroupID(part *multipart.Part) error {
	var buf bytes.Buffer

	if _, err := io.Copy(&buf, part); err != nil {
		return fmt.Errorf("reading group_id: %w", err)
	}

	f.item.GroupID = strings.Split(buf.String(), ",")

	return nil
}

func (f *documentUploadForm) parseDocument(part *multipart.Part) error {
	file, err := mapPartToFile(part)
	if err != nil {
		return fmt.Errorf("mapping document file: %w", err)
	}

	f.item.Source = *file

	return nil
}

func (f *documentUploadForm) Parse(reader *multipart.Reader) (model.Document, error) {
	for {
		part, err := reader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}

			return model.Document{}, fmt.Errorf("reading multipart part: %w", err)
		}

		if cb, exists := f.callbacks[part.FormName()]; exists {
			if err := cb(part); err != nil {
				return model.Document{}, err
			}
		}
	}

	return f.item, nil
}

func mapUploadRequestToCommand(
	in openapi.PostDocumentUploadRequestObject,
) (model.DocumentSaveCommand, error) {
	uploadForm := newDocumentUploadForm()

	doc, err := uploadForm.Parse(in.Body)
	if err != nil {
		return model.DocumentSaveCommand{}, fmt.Errorf("parse form: %w", err)
	}

	doc.Name = doc.Source.Name

	return model.DocumentSaveCommand{
		Item: doc,
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

func mapDocumentFileRequestToQuery(
	in openapi.GetDocumentIdDownloadRequestObject,
) model.DocumentQuery {
	return model.DocumentQuery{
		ID:          in.DocumentId,
		WithContent: true,
		Include: []model.DocumentQueryInclude{
			model.DocumentQueryIncludeSource,
		},
	}
}

func mapFileToDownloadResponse(
	in model.File,
) openapi.GetDocumentIdDownload200ApplicationoctetStreamResponse {
	return openapi.GetDocumentIdDownload200ApplicationoctetStreamResponse{
		Body: bytes.NewReader(in.Content),
		Headers: openapi.GetDocumentIdDownload200ResponseHeaders{
			ContentDisposition: fmt.Sprintf(`attachment; filename="%s"`, in.Name),
		},
		ContentLength: int64(len(in.Content)),
	}
}

func mapPartToFile(in *multipart.Part) (*model.File, error) {
	if in == nil {
		return nil, errNilPart
	}

	buf := &bytes.Buffer{}
	hasher := sha256.New()

	_, err := io.Copy(io.MultiWriter(buf, hasher), in)
	if err != nil {
		return nil, fmt.Errorf("read part: %w", err)
	}

	return &model.File{
		Name:        in.FileName(),
		Content:     buf.Bytes(),
		Sha256:      hex.EncodeToString(hasher.Sum(nil)),
		LastUpdated: time.Time{},
	}, nil
}

func valOrDefault[T any](in *T) T {
	if in == nil {
		var t T
		return t
	}

	return *in
}
