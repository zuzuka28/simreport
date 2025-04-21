package elasticstorage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
)

type Config struct {
	Index string `yaml:"index"`
}

type Storage struct {
	opts Config
	r    *elasticsearch.Client
}

func New(opts Config, cli *elasticsearch.Client) *Storage {
	return &Storage{
		opts: opts,
		r:    cli,
	}
}

func (s *Storage) Fetch(ctx context.Context, key string) ([]string, error) {
	resp, err := s.r.Get(
		s.opts.Index,
		key,
		s.r.Get.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("fetch from elasticsearch: %w", err)
	}

	if resp.IsError() {
		if resp.StatusCode == http.StatusNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("%w: %d", errBadStatusCode, resp.StatusCode)
	}

	raw := new(fetchResponse)

	if err := json.NewDecoder(resp.Body).Decode(raw); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return raw.Source.Content, nil
}

func (s *Storage) Insert(ctx context.Context, key string, value string) error {
	prev, err := s.Fetch(ctx, key)
	if err != nil {
		return fmt.Errorf("fetch: %w", err)
	}

	vals := addUniqValue(prev, value)

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(&storedDocumentSource{
		ID:      key,
		Content: vals,
	}); err != nil {
		return fmt.Errorf("encode content: %w", err)
	}

	resp, err := s.r.Index(
		s.opts.Index,
		&buf,
		s.r.Index.WithContext(ctx),
		s.r.Index.WithDocumentID(key),
	)
	if err != nil {
		return fmt.Errorf("set to elasticsearch: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("%w: %d", errBadStatusCode, resp.StatusCode)
	}

	return nil
}

func addUniqValue(to []string, val string) []string {
	to = append(to, val)

	uniq := make(set)

	for _, v := range to {
		uniq[v] = struct{}{}
	}

	res := make([]string, 0, len(uniq))

	for _, v := range to {
		_, ok := uniq[v]
		if ok {
			res = append(res, v)
		}

		delete(uniq, v)
	}

	return res
}
