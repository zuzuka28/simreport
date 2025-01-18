package redisstorage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	redis "github.com/redis/go-redis/v9"
)

type Config struct {
	DSN string `yaml:"dsn"`
}

type Storage struct {
	r *redis.Client
}

func New(cli *redis.Client) *Storage {
	return &Storage{
		r: cli,
	}
}

func (s *Storage) Fetch(ctx context.Context, key string) ([]string, error) {
	resp, err := s.r.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, fmt.Errorf("fetch from redis: %w", err)
	}

	raw := make(set)

	if err := json.Unmarshal([]byte(resp), &raw); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	result := make([]string, 0, len(raw))

	for k := range raw {
		result = append(result, k)
	}

	return result, nil
}

func (s *Storage) Insert(ctx context.Context, key string, value string) error {
	resp, err := s.r.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return fmt.Errorf("fetch from redis: %w", err)
	}

	raw := make(set)

	if resp != "" {
		if err := json.Unmarshal([]byte(resp), &raw); err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}
	}

	raw[value] = struct{}{}

	m, err := json.Marshal(raw)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	err = s.r.Set(ctx, key, m, 0).Err()
	if err != nil {
		return fmt.Errorf("set to redis: %w", err)
	}

	return nil
}
