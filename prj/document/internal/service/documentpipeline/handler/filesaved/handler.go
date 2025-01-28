package filesaved

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/document/internal/model"
)

type Handler struct {
	ds DocumentService
}

func NewHandler(
	ds DocumentService,
) *Handler {
	return &Handler{
		ds: ds,
	}
}

func (h *Handler) Serve(ctx context.Context, documentID string) error {
	doc, err := h.ds.Fetch(ctx, model.DocumentQuery{
		ID:          documentID,
		WithContent: true,
		Include: []model.DocumentQueryInclude{
			model.DocumentQueryIncludeSource,
		},
	})
	if err != nil {
		return fmt.Errorf("fetch file: %w", err)
	}

	parsed, err := h.ds.Parse(ctx, doc.Source)
	if err != nil {
		return fmt.Errorf("parse document: %w", err)
	}

	if _, err := h.ds.Save(ctx, model.DocumentSaveCommand{
		Item: mergeDocuments(doc, parsed),
	}); err != nil {
		return fmt.Errorf("save parsed document: %w", err)
	}

	return nil
}

func mergeDocuments(main, delt model.Document) model.Document {
	return model.Document{
		ParentID:    main.ParentID,
		Name:        main.Name,
		LastUpdated: main.LastUpdated,
		Version:     main.Version,
		GroupID:     main.GroupID,
		SourceID:    main.SourceID,
		TextID:      delt.TextID,
		ImageIDs:    delt.ImageIDs,
		WithContent: false,
		Source:      main.Source,
		Text:        delt.Text,
		Images:      delt.Images,
	}
}
