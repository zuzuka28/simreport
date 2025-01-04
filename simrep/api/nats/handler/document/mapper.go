package document

import "simrep/internal/model"

func mapDocumentToResponse(in model.Document) Document {
	return Document{
		ID:   in.ID,
		Text: in.Text.Content,
	}
}
