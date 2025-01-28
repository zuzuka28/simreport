package shingleindex

import (
	"context"
	"shingleindex/internal/model"

	"github.com/nats-io/nats.go/micro"
)

func (h *Handler) SearchSimilar(msg micro.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	id := string(msg.Data())

	doc, err := h.ds.Fetch(ctx, model.DocumentQuery{
		ID: id,
	})
	if err != nil {
		_ = msg.Error(mapErrorToStatus(err), err.Error(), nil)
		return
	}

	res, err := h.s.SearchSimilar(ctx, model.DocumentSimilarQuery{
		ID:   id,
		Item: doc,
	})
	if err != nil {
		_ = msg.Error(mapErrorToStatus(err), err.Error(), nil)
		return
	}

	_ = msg.RespondJSON(mapDocumentToResponse(res))
}
