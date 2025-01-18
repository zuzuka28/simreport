package shingleindex

import (
	"context"
	"shingleindex/internal/model"
	"time"

	"github.com/nats-io/nats.go/micro"
)

const requestTimeout = 60 * time.Second

type Handler struct {
	s  Service
	ds DocumentService
}

func NewHandler(
	s Service,
	ds DocumentService,
) *Handler {
	return &Handler{
		s:  s,
		ds: ds,
	}
}

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
