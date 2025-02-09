package shingleindex

import (
	"hash/fnv"
)

const (
	shingleSize  = 4
	permutations = 512
	seed         = 42
)

//nolint:gochecknoglobals
var hasher = func(b []byte) uint64 {
	h := fnv.New64a()
	_, _ = h.Write(b)
	return h.Sum64()
}
