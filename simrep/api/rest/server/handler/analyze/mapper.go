package analyze

import (
	"bytes"
	"context"
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

func makeMapSearchSimilarRequestToQuery(
	s Service,
	dp DocumentParser,
) func(
	ctx context.Context,
	in openapi.PostAnalyzeSimilarRequestObject,
) (model.AnalyzedDocumentSimilarQuery, error) {
	return func(
		ctx context.Context,
		in openapi.PostAnalyzeSimilarRequestObject,
	) (model.AnalyzedDocumentSimilarQuery, error) {
		file, err := fileFromMultipart(in.Body, "document")
		if err != nil {
			return model.AnalyzedDocumentSimilarQuery{},
				fmt.Errorf("retrieve file from multipart: %w", err)
		}

		parsed, err := dp.Parse(ctx, file)
		if err != nil {
			return model.AnalyzedDocumentSimilarQuery{},
				fmt.Errorf("parse input file: %w", err)
		}

		analyzed, err := s.Analyze(ctx, parsed)
		if err != nil {
			return model.AnalyzedDocumentSimilarQuery{},
				fmt.Errorf("analyze input file: %w", err)
		}

		return model.AnalyzedDocumentSimilarQuery{
			Item: analyzed,
		}, nil
	}
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

func sha256String(in []byte) string {
	hash := sha256.New()
	_, _ = hash.Write(in)

	return hex.EncodeToString(hash.Sum(nil))
}

func mapMatchesToSearchSimilarResponse(in []model.AnalyzedDocumentMatch) openapi.PostAnalyzeSimilar200JSONResponse {
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

	return openapi.PostAnalyzeSimilar200JSONResponse{
		SimilaritySearchResultJSONResponse: openapi.SimilaritySearchResultJSONResponse{
			Documents: &docs,
		},
	}
}
