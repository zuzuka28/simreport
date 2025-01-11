package document

import (
	"context"
	"document/internal/model"
	"time"

	"github.com/nats-io/nats.go/micro"
)

const requestTimeout = 60 * time.Second

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

func (h *Handler) Fetch(msg micro.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	id := string(msg.Data())

	doc, err := h.s.Fetch(ctx, model.DocumentQuery{
		ID:          id,
		WithContent: true,
		Include: []model.DocumentQueryInclude{
			model.DocumentQueryIncludeSource,
			model.DocumentQueryIncludeText,
			model.DocumentQueryIncludeImages,
		},
	})
	if err != nil {
		_ = msg.Error(mapErrorToStatus(err), err.Error(), nil)
		return
	}

	_ = msg.RespondJSON(mapDocumentToResponse(doc))
}
