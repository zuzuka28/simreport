package semanticindexclient

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/zuzuka28/simreport/prj/document/internal/model"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

var errInternal = errors.New("internal error")

func parseSearchSimilarResponse(in []byte) ([]*model.SimilarityMatch, error) {
	if len(in) == 0 {
		return nil, nil
	}

	var raw []*documentSimilarMatch

	if err := json.Unmarshal(in, &raw); err != nil {
		return nil, fmt.Errorf("unmarshal raw: %w", err)
	}

	items := make([]*model.SimilarityMatch, 0, len(raw))
	for _, v := range raw {
		items = append(items, mapDocumentSimilarMatchToModel(v))
	}

	return items, nil
}

func mapDocumentSimilarMatchToModel(in *documentSimilarMatch) *model.SimilarityMatch {
	if in == nil {
		return nil
	}

	return &model.SimilarityMatch{
		ID:            in.ID,
		Rate:          in.Rate,
		Highlights:    in.Highlights,
		SimilarImages: in.SimilarImages,
	}
}

func isErr(in *nats.Msg) error {
	status := in.Header.Get(micro.ErrorCodeHeader)
	if status == "" {
		return nil
	}

	switch status {
	case "404":
		return model.ErrNotFound

	case "400":
		return model.ErrInvalid

	default:
		return errInternal
	}
}
