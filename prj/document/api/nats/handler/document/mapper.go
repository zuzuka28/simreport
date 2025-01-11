package document

import (
	"errors"
	"document/internal/model"
)

func mapDocumentToResponse(in model.Document) Document {
	return Document{
		ID:   in.ID,
		Text: in.Text.Content,
	}
}

func mapErrorToStatus(err error) string {
	if errors.Is(err, model.ErrNotFound) {
		return "404"
	}

	if errors.Is(err, model.ErrInvalid) {
		return "400"
	}

	return "500"
}
