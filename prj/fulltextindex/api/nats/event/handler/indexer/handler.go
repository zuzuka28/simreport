package indexer

import (
	"context"
	"fulltextindex/internal/model"
	"time"

	"github.com/nats-io/nats.go"
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

func (h *Handler) Save(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	id := string(msg.Data)

	doc, err := h.ds.Fetch(ctx, model.DocumentQuery{
		ID: id,
	})
	if err != nil {
		_ = msg.Nak()
		return
	}

	err = h.s.Save(ctx, model.DocumentSaveCommand{
		Item: doc,
	})
	if err != nil {
		_ = msg.Nak()
		return
	}

	_ = msg.AckSync()
}
