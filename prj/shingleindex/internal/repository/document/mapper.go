package document

import (
	"encoding/json"
	"errors"
	"fmt"
	"shingleindex/internal/model"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

var errInternal = errors.New("internal error")

func parseFetchDocumentResponse(in []byte) (model.Document, error) {
	if len(in) == 0 {
		return model.Document{}, nil
	}

	var raw document

	if err := json.Unmarshal(in, &raw); err != nil {
		return model.Document{}, fmt.Errorf("unmarshal raw: %w", err)
	}

	return model.Document{
		ID:   raw.ID,
		Text: raw.Text,
	}, nil
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
