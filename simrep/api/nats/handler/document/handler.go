package document

import (
	"context"
	"encoding/json"
	"simrep/internal/model"
	"time"

	"github.com/nats-io/nats.go"
)

const requestTimeout = 5 * time.Second

type Handler struct {
	s Service
}

func NewHandler(
	s Service,
) *Handler {
	return &Handler{
		s: s,
	}
}

func (h *Handler) Fetch(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	id := string(msg.Data)

	doc, err := h.s.Fetch(ctx, model.DocumentQuery{
		ID:          id,
		WithContent: true,
	})
	if err != nil {
		return
	}

	res, err := json.Marshal(mapDocumentToResponse(doc))
	if err != nil {
		return
	}

	_ = msg.Respond(res)
}
