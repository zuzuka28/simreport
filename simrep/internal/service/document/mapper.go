package document

import "simrep/internal/model"

func mapDocumentWithContentToDocument(in model.Document) model.Document {
	return model.Document{
		ID:          in.ID,
		Name:        in.Name,
		LastUpdated: in.Source.LastUpdated,
		ImageIDs:    in.ImageIDs,
		TextID:      in.TextID,
		WithContent: false,
		Source:      model.File{}, //nolint:exhaustruct
		Text:        model.File{}, //nolint:exhaustruct
		Images:      []model.File{},
	}
}
