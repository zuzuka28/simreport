package shingleindex

import (
	"errors"
	"shingleindex/internal/model"
)

func mapDocumentToResponse(in []*model.DocumentSimilarMatch) []documentSimilarMatch {
	items := make([]documentSimilarMatch, 0, len(in))
	for _, v := range in {
		items = append(items, documentSimilarMatch{
			ID:            v.ID,
			Rate:          v.Rate,
			Highlights:    v.Highlights,
			SimilarImages: v.SimilarImages,
		})
	}
	return items
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
