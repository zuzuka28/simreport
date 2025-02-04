package vectorizer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zuzuka28/simreport/prj/semanticindex/internal/model"
	client "github.com/zuzuka28/simreport/prj/semanticindex/internal/repository/vectorizer/vectorizerclient"
)

func (s *Repository) TextToVector(
	ctx context.Context,
	params model.VectorizeTextParams,
) (model.Vector, error) {
	requestBody := client.VectorizeTextVectorizeTextPostJSONRequestBody{
		Text: params.Text,
	}

	resp, err := s.cli.VectorizeTextVectorizeTextPostWithResponse(
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
