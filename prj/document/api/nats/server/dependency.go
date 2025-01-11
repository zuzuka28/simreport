package server

import (
	"github.com/nats-io/nats.go/micro"
)

type (
	DocumentHandler interface {
		Fetch(req micro.Request)
	}
)
