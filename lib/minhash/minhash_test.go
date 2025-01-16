package minhash_test

import (
	"hash/fnv"
	"testing"

	"github.com/zuzuka28/simreport/lib/minhash"
)

func jaccardSimilarity[T comparable](l1, l2 []T) float64 {
	s1 := make(map[T]struct{}, len(l1))

	for _, v := range l1 {
		s1[v] = struct{}{}
	}

	s2 := make(map[T]struct{}, len(l2))

	for _, v := range l2 {
		s2[v] = struct{}{}
	}

	var intersection []T

	for k := range s1 {
		if _, ok := s2[k]; ok {
			intersection = append(intersection, k)
		}
	}

	union := make(map[T]struct{})

	for k := range s1 {
		union[k] = struct{}{}
	}

	for k := range s2 {
		union[k] = struct{}{}
	}

	return float64(len(intersection)) / float64(len(union))
}

func TestMinhash(t *testing.T) {
	t.Parallel()

	var (
		permutationsNum = 512
		hasher          = func(b []byte) uint64 {
			h := fnv.New64a()
			h.Write(b)
			return h.Sum64()
		}
		seed = 1
	)

	tests := []struct {
		s1 []string
		s2 []string
	}{
		{
			[]string{"hello"},
			[]string{"hello"},
		},
		{
			[]string{
				"minhash",
				"is",
				"a",
				"probabilistic",
				"data",
				"structure",
				"for",
				"estimating",
				"the",
				"similarity",
				"between",
				"datasets",
			},
			[]string{
				"minhash",
				"is",
				"a",
				"probability",
				"data",
				"structure",
				"for",
				"estimating",
				"the",
				"similarity",
				"between",
				"documents",
			},
		},
		{
			[]string{"hello", "world"},
			[]string{"hello"},
		},
		{
			[]string{"hello", "world", "foo", "baz", "bar", "zomg"},
			[]string{"goodbye", "world", "foo", "qux", "bar", "zomg"},
		},
	}

	for _, tt := range tests {
		m1 := minhash.New(permutationsNum, hasher, seed)
		m2 := minhash.New(permutationsNum, hasher, seed)

		for _, s := range tt.s1 {
			m1.Push([]byte(s))
		}

		for _, s := range tt.s2 {
			m2.Push([]byte(s))
		}

		sim, _ := m1.Similarity(m2)
		t.Log("similarity by minhash", sim)

		t.Log("similarity by jaccard", jaccardSimilarity(tt.s1, tt.s2))
	}
}
