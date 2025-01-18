package shingleindex

import (
	"hash/fnv"
	"shingleindex/internal/model"
)

const (
	shingleSize  = 4
	permutations = 512
	seed         = 42
)

type documentMatch struct {
	*model.DocumentSimilarMatch
	*model.MinhashSimilarMatch
	text     string
	shingles map[string]struct{}
}

var hasher = func(b []byte) uint64 {
	h := fnv.New64a()
	h.Write(b)
	return h.Sum64()
}
