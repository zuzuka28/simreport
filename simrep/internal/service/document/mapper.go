package document

import "simrep/internal/model"

func mapDocumentWithContentToDocument(in model.Document) model.Document {
	return model.Document{
		ID:          in.ID,
		Name:        in.Name,
		ImageIDs:    in.ImageIDs,
		TextContent: in.TextContent,
		LastUpdated: in.Source.LastUpdated,
		WithContent: false,
		Source:      model.File{},
		Images:      []model.File{},
	}
}
