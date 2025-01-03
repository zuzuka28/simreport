package filesaved

import (
	"context"
	"fmt"
	"simrep/internal/model"
)

type Handler struct {
	fs FileService
	ds DocumentService
}

func NewHandler(
	fs FileService,
	ds DocumentService,
) *Handler {
	return &Handler{
		fs: fs,
		ds: ds,
	}
}

func (h *Handler) Serve(ctx context.Context, documentID string) error {
	doc, err := h.fs.Fetch(ctx, model.FileQuery{
		ID: documentID,
	})
	if err != nil {
		return fmt.Errorf("fetch file: %w", err)
	}

	parsed, err := h.ds.Parse(ctx, doc)
	if err != nil {
		return fmt.Errorf("parse document: %w", err)
	}

	if err := h.ds.Save(ctx, model.DocumentSaveCommand{
		Item: parsed,
	}); err != nil {
		return fmt.Errorf("save parsed document: %w", err)
	}

	return nil
}
