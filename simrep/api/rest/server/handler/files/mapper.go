package files

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

var errNilPart = errors.New("nil part")

func mapUploadFilesToCommand(
	in openapi.PostFilesUploadRequestObject,
) (model.DocumentFileUploadManyCommand, error) {
	var items []model.File

	for {
		part, err := in.Body.NextPart()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return model.DocumentFileUploadManyCommand{}, fmt.Errorf("read part: %w", err)
		}

		if part.FormName() != "files" {
			continue
		}

		item, err := mapPartToFile(part)
		if err != nil {
			return model.DocumentFileUploadManyCommand{}, fmt.Errorf("map part to file: %w", err)
		}

		items = append(items, *item)
	}

	return model.DocumentFileUploadManyCommand{
		Items: items,
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
		Content:     buf.Bytes(),
		Sha256:      hex.EncodeToString(sha256.New().Sum(buf.Bytes())),
		LastUpdated: time.Time{},
	}, nil
}

func mapCommandToUploadSuccess(cmd model.DocumentFileUploadManyCommand) openapi.UploadSuccess {
	ids := make([]string, 0, len(cmd.Items))

	for _, v := range cmd.Items {
		ids = append(ids, v.Sha256)
	}

	return openapi.UploadSuccess{
		FileIds: &ids,
	}
}
