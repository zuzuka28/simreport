package indexer

import (
	"context"

	"github.com/zuzuka28/simreport/prj/semanticindex/internal/model"

	"github.com/nats-io/nats.go"
)

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
