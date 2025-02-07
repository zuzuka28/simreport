package documentstatus

import (
	"github.com/nats-io/nats.go/jetstream"
)

const subject = "documentstatus"

type Repository struct {
	kv jetstream.KeyValue
	p  jetstream.Publisher

	m Metrics
}

func NewRepository(
	kv jetstream.KeyValue,
	p jetstream.Publisher,
	m Metrics,
) *Repository {
	return &Repository{
		kv: kv,
		p:  p,
		m:  m,
	}
}
