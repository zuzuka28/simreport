package minhashlsh

import (
	"golang.org/x/net/context"
)

type Storage interface {
	Fetch(ctx context.Context, key string) ([]string, error)
	Insert(ctx context.Context, key string, value string) error
}
