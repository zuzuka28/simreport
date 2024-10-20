package documentsaved

import (
	"context"
	"fmt"
	"simrep/internal/model"
)

type Handler struct {
	ds DocumentService
	as AnalyzeService
}

func NewHandler(
	ds DocumentService,
	as AnalyzeService,
) *Handler {
	return &Handler{
		ds: ds,
		as: as,
	}
}

func (h *Handler) Serve(ctx context.Context, documentID string, _ any) error {
	doc, err := h.ds.Fetch(ctx, model.DocumentQuery{
		ID:          documentID,
		WithContent: true,
	})
	if err != nil {
		return fmt.Errorf("fetch document: %w", err)
	}

	analyzed, err := h.as.Analyze(ctx, doc)
	if err != nil {
		return fmt.Errorf("analyze document: %w", err)
	}

	if err := h.as.Save(ctx, model.AnalyzedDocumentSaveCommand{
		Item: analyzed,
	}); err != nil {
		return fmt.Errorf("save analyzed document: %w", err)
	}

	return nil
}
