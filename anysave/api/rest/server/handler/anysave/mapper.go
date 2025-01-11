package anysave

import (
	openapi "anysave/api/rest/gen"
	"anysave/internal/model"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"time"
)

var (
	errNilPart    = errors.New("nil part")
	errNoDocument = errors.New("no document")
)

func mapUploadRequestToCommand(
	in openapi.PostUploadRequestObject,
) (model.FileSaveCommand, error) {
	file, err := fileFromMultipart(in.Body, "document")
	if err != nil {
		return model.FileSaveCommand{},
			fmt.Errorf("retrieve file from multipart: %w", err)
	}

	return model.FileSaveCommand{
		Bucket: "",
		Item:   file,
	}, nil
}

func fileFromMultipart(r *multipart.Reader, name string) (model.File, error) {
	var (
		file     model.File
		hasInput bool
	)

	for {
		part, err := r.NextPart()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return model.File{}, fmt.Errorf("read part: %w", err)
		}

		if part.FormName() != name {
			continue
		}

		item, err := mapPartToFile(part)
		if err != nil {
			return model.File{}, fmt.Errorf("map part to file: %w", err)
		}

		file = *item
		hasInput = true

		break
	}

	if !hasInput {
		return model.File{}, errNoDocument
	}

	return file, nil
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
	cmd model.FileSaveCommand,
) openapi.PostUpload200JSONResponse {
	return openapi.PostUpload200JSONResponse{
		UploadSuccessJSONResponse: openapi.UploadSuccessJSONResponse(
			openapi.UploadSuccess{
				DocumentID: &cmd.Item.Sha256,
			},
		),
	}
}

func mapDocumentFileRequestToQuery(
	in openapi.GetDocumentIdDownloadRequestObject,
) model.FileQuery {
	return model.FileQuery{
		ID: in.DocumentId,
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

func sha256String(in []byte) string {
	hash := sha256.New()
	_, _ = hash.Write(in)

	return hex.EncodeToString(hash.Sum(nil))
}
