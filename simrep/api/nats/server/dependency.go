package server

import "github.com/nats-io/nats.go"

type (
	DocumentHandler interface {
		Fetch(msg *nats.Msg)
	}
)
