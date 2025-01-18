package minhashlsh

import (
	"context"
	"strconv"
)

type table struct {
	prefix string
	s      Storage
}

func newTable(prefix string, s Storage) *table {
	return &table{
		prefix: prefix,
		s:      s,
	}
}

func (t *table) Fetch(ctx context.Context, band uint64) ([]string, error) {
	key := t.prefix + ":" + strconv.FormatUint(band, 10)

	return t.s.Fetch(ctx, key)
}

func (t *table) Insert(ctx context.Context, band uint64, val string) error {
	key := t.prefix + ":" + strconv.FormatUint(band, 10)

	return t.s.Insert(ctx, key, val)
}
