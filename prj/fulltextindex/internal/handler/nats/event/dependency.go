package server

import "github.com/nats-io/nats.go"

type (
	IndexerHandler interface {
		Save(msg *nats.Msg)
	}
)
