package similarity

import (
	"hash/fnv"

	"github.com/zuzuka28/simreport/prj/similarity/internal/model"
)

const (
	shingleSize  = 4
	permutations = 512
	seed         = 42
)

type match struct {
	*model.SimilarityMatch

	docs []model.Document
	textid   string

	text     string
	shingles map[string]struct{}
}

var hasher = func(b []byte) uint64 {
	h := fnv.New64a()
	h.Write(b)
	return h.Sum64()
}
