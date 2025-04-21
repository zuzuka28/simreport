package shingleindex

import "hash/fnv"

const (
	prefix       = "shingleindex"
	permutations = 512
	bands        = 64
)

//nolint:gochecknoglobals
var hasher = func(b []byte) uint64 {
	h := fnv.New64a()
	h.Write(b)
	return h.Sum64()
}
