package server

import (
	"github.com/nats-io/nats.go/micro"
)

type (
	FileindexHandler interface {
		SearchSimilar(msg micro.Request)
	}
)
