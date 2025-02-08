package vectorizer

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/zuzuka28/simreport/prj/semanticindex/internal/model"
	client "github.com/zuzuka28/simreport/prj/semanticindex/internal/repository/vectorizer/vectorizerclient"
)

func (s *Repository) TextToVector(
	ctx context.Context,
	params model.VectorizeTextParams,
) (model.Vector, error) {
	const op = "textToVector"

	requestBody := client.VectorizeTextVectorizeTextPostJSONRequestBody{
		Text: params.Text,
	}

	t := time.Now()

	resp, err := s.cli.VectorizeTextVectorizeTextPostWithResponse(
		ctx,
		requestBody,
	)
	if err != nil {
		s.m.IncVectorizerRequests(op, metricsError, time.Since(t).Seconds())
		return model.Vector{}, fmt.Errorf("vectorize text: %w", err)
	}

	s.m.IncVectorizerRequests(op, strconv.Itoa(resp.StatusCode()), time.Since(t).Seconds())

	if resp.StatusCode() != http.StatusOK {
		return model.Vector{}, fmt.Errorf("%w: status %d", errBadResponse, resp.StatusCode())
	}

	return convertVector(resp.JSON200.Vector), nil
}
