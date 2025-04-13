package indexer

import (
	"context"

	"github.com/zuzuka28/simreport/prj/fulltextindex/internal/model"

	"github.com/nats-io/nats.go"
)

const bucketTexts = "texts"

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

	textfile, err := h.fs.Fetch(ctx, model.FileQuery{
		Bucket: bucketTexts,
		ID:     doc.TextID,
	})
	if err != nil {
		_ = msg.Nak()
		return
	}

	doc.Text = textfile.Content

	err = h.s.Save(ctx, model.DocumentSaveCommand{
		Item: doc,
	})
	if err != nil {
		_ = msg.Nak()
		return
	}

	_ = msg.AckSync()
}
