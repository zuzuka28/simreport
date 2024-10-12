package docprepsrv

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"simrep/internal/model"
	client "simrep/pkg/docprepsrv"
)

func mapDocToPreprocessRequest(in []byte) (io.Reader, string, error) {
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)

	part, err := writer.CreateFormFile("file", "file.docx")
	if err != nil {
		return nil, "", fmt.Errorf("create form field: %w", err)
	}

	_, err = io.Copy(part, bytes.NewReader(in))
	if err != nil {
		return nil, "", fmt.Errorf("copy data to multipart: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("close multipart: %w", err)
	}

	return buf, writer.FormDataContentType(), nil
}

func mapDocResponseToModel(
	in *client.DocumentResponse,
) *model.DocumentFull {
	if in == nil {
		return nil
	}

	rawimages := zeroOrValue(in.Images)
	images := make([]*model.Image, 0, len(rawimages))

	for _, v := range rawimages {
		images = append(images, mapImageToModel(&v))
	}

	return &model.DocumentFull{
		ID:              in.Id,
		Images:          images,
		SbertTextVector: zeroOrValue(in.SbertTextVector),
		Sha256:          in.Sha256,
		SourceBytes:     in.SourceBytes,
		TextContent:     in.TextContent,
	}
}

func mapImageToModel(
	in *client.ImageResponse,
) *model.Image {
	if in == nil {
		return nil
	}

	return &model.Image{
		ClipImageVector: zeroOrValue(in.ClipImageVector),
		Fname:           in.Fname,
		Hashes:          mapImageHashesToModel(in.Hashes),
		Sha256:          in.Sha256,
		SourceBytes:     in.SourceBytes,
	}
}

func mapImageHashesToModel(
	in *client.ImageHashesResponse,
) *model.ImageHashes {
	if in == nil {
		return nil
	}

	return &model.ImageHashes{
		Ahash:       in.Ahash,
		AhashVector: in.AhashVector,
		Dhash:       in.Dhash,
		DhashVector: in.DhashVector,
		Phash:       in.Phash,
		PhashVector: in.PhashVector,
		Whash:       in.Whash,
		WhashVector: in.WhashVector,
	}
}

func zeroOrValue[T any](in *T) T {
	if in == nil {
		var t T
		return t
	}

	return *in
}
