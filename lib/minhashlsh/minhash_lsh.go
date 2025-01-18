package minhashlsh

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/zuzuka28/simreport/lib/minhash"
)

type HashFunc func([]byte) uint64

type MinhashLSH struct {
	h HashFunc

	bandSize int
	tables   []*table
}

func New(
	s Storage,
	prefix string,
	permutations int,
	bands int,
	h HashFunc,
) *MinhashLSH {
	tables := make([]*table, bands)

	for i := range bands {
		p := prefix + "_band_" + strconv.Itoa(i)
		tables[i] = newTable(p, s)
	}

	return &MinhashLSH{
		h:        h,
		bandSize: permutations / bands,
		tables:   tables,
	}
}

func (l *MinhashLSH) Insert(
	ctx context.Context,
	key string,
	mh *minhash.MinHash,
) error {
	mhv := mh.Values()

	for i, tab := range l.tables {
		bandKey := l.h(uintsToBytes(mhv[i*l.bandSize : i*l.bandSize+l.bandSize]))

		if err := tab.Insert(ctx, bandKey, key); err != nil {
			return fmt.Errorf("insert: %w", err)
		}
	}

	return nil
}

func (l *MinhashLSH) Query(
	ctx context.Context,
	mh *minhash.MinHash,
) ([]string, error) {
	candidates := make(map[string]struct{})

	mhv := mh.Values()

	for i, tab := range l.tables {
		bandKey := l.h(uintsToBytes(mhv[i*l.bandSize : i*l.bandSize+l.bandSize]))

		c, err := tab.Fetch(ctx, bandKey)
		if err != nil {
			return nil, fmt.Errorf("fetch: %w", err)
		}

		for _, v := range c {
			candidates[v] = struct{}{}
		}
	}

	result := make([]string, 0, len(candidates))

	for k := range candidates {
		result = append(result, k)
	}

	return result, nil
}

func uintsToBytes(in []uint64) []byte {
	var buf bytes.Buffer

	for _, v := range in {
		_ = binary.Write(&buf, binary.LittleEndian, v)
	}

	return buf.Bytes()
}
