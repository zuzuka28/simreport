package vectorizer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"simrep/internal/model"
	client "simrep/pkg/vectorizerclient"
)

var errBadResponse = errors.New("bad response")

type Service struct {
	client client.ClientWithResponsesInterface
}

func NewService(
	cli client.ClientWithResponsesInterface,
) *Service {
	return &Service{
		client: cli,
	}
}

func (s *Service) TextToVector(
	ctx context.Context,
	params model.VectorizeTextParams,
) (model.Vector, error) {
	requestBody := client.VectorizeTextVectorizeTextPostJSONRequestBody{
		Text: params.Text,
	}

	resp, err := s.client.VectorizeTextVectorizeTextPostWithResponse(
		ctx,
		requestBody,
	)
	if err != nil {
		return model.Vector{}, fmt.Errorf("vectorize text: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return model.Vector{}, fmt.Errorf("%w: status %d", errBadResponse, resp.StatusCode())
	}

	return convertVector(resp.JSON200.Vector), nil
}

func (s *Service) ImageToVector(
	ctx context.Context,
	params model.VectorizeImageParams,
) (model.Vector, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	filePart, err := writer.CreateFormFile("file", params.Image.Name)
	if err != nil {
		return model.Vector{}, fmt.Errorf("create formfile: %w", err)
	}

	_, err = io.Copy(filePart, bytes.NewReader(params.Image.Content))
	if err != nil {
		return model.Vector{}, fmt.Errorf("write filepart: %w", err)
	}

	_ = writer.Close()

	resp, err := s.client.VectorizeImageVectorizeImagePostWithBodyWithResponse(
		ctx,
		writer.FormDataContentType(),
		body,
	)
	if err != nil {
		return model.Vector{}, fmt.Errorf("vectorize image: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return model.Vector{}, fmt.Errorf("%w: status %d", errBadResponse, resp.StatusCode())
	}

	return convertVector(resp.JSON200.Vector), nil
}

func (s *Service) ImageToHashes(
	ctx context.Context,
	params model.HashImageParams,
) (model.HashImage, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	filePart, err := writer.CreateFormFile("file", params.Image.Name)
	if err != nil {
		return model.HashImage{}, fmt.Errorf("create formfile: %w", err)
	}

	_, err = io.Copy(filePart, bytes.NewReader(params.Image.Content))
	if err != nil {
		return model.HashImage{}, fmt.Errorf("write filepart: %w", err)
	}

	_ = writer.Close()

	resp, err := s.client.HashImageHashImagePostWithBodyWithResponse(
		ctx,
		writer.FormDataContentType(),
		body,
	)
	if err != nil {
		return model.HashImage{}, fmt.Errorf("hash image: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return model.HashImage{}, fmt.Errorf("%w: status %d", errBadResponse, resp.StatusCode())
	}

	return model.HashImage{
		Ahash:       resp.JSON200.Ahash,
		Dhash:       resp.JSON200.Dhash,
		Phash:       resp.JSON200.Phash,
		Whash:       resp.JSON200.Whash,
		AhashVector: convertVector(resp.JSON200.AhashVector),
		DhashVector: convertVector(resp.JSON200.DhashVector),
		PhashVector: convertVector(resp.JSON200.PhashVector),
		WhashVector: convertVector(resp.JSON200.WhashVector),
	}, nil
}

func convertVector(f32 []float32) []float64 {
	f64 := make([]float64, len(f32))
	for i, v := range f32 {
		f64[i] = float64(v)
	}

	return f64
}
