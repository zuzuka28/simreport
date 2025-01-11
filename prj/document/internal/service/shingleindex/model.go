package shingleindex

import "document/internal/model"

type documentMatch struct {
	*model.DocumentSimilarMatch
	text     string
	shingles map[string]struct{}
}
