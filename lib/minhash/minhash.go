package minhash

import (
	"errors"
	"math/rand"
)

var ErrDifferentPermutationsNum = errors.New("different permutations num")

var (
	mersennePrime uint64 = ((1 << 61) - 1)
	maxHash       uint64 = ((1 << 32) - 1)
)

type HashFunc func([]byte) uint64

type MinHash struct {
	numPermutations int
	seed            int

	h            HashFunc
	permutations [][2]uint64
	values       []uint64
}

func New(
	permutations int,
	hash HashFunc,
	seed int,
) *MinHash {
	return &MinHash{
		numPermutations: permutations,
		seed:            seed,
		h:               hash,
		permutations:    genPermutations(permutations, seed),
		values:          genValues(permutations),
	}
}

func (m *MinHash) Push(v []byte) {
	hv := m.h(v)
	for i, perm := range m.permutations {
		m.values[i] = min(((hv*perm[0]+perm[1])%mersennePrime)&maxHash, m.values[i])
	}
}

func (m *MinHash) Values() []uint64 {
	dst := make([]uint64, len(m.values))
	copy(dst, m.values)

	return dst
}

func (m *MinHash) Merge(other *MinHash) error {
	if m.numPermutations != other.numPermutations {
		return ErrDifferentPermutationsNum
	}

	for i, v := range other.values {
		if v < m.values[i] {
			m.values[i] = v
		}
	}
	return nil
}

func (m *MinHash) Similarity(other *MinHash) (float64, error) {
	if m.numPermutations != other.numPermutations {
		return 0, ErrDifferentPermutationsNum
	}

	count := 0
	for i, v := range m.values {
		if v == other.values[i] {
			count++
		}
	}

	return float64(count) / float64(m.numPermutations), nil
}

func genPermutations(count, seed int) [][2]uint64 {
	rnd := rand.New(rand.NewSource(int64(seed)))

	vals := make([][2]uint64, count)
	for i := range vals {
		vals[i] = [2]uint64{
			1 + uint64(rnd.Uint64()%(mersennePrime-1)),
			0 + uint64(rnd.Uint64()%(mersennePrime-0)),
		}
	}

	return vals
}

func genValues(size int) []uint64 {
	vals := make([]uint64, size)
	for i := range vals {
		vals[i] = maxHash
	}

	return vals
}
