package consumer

import (
	"context"
	"simrep/pkg/rabbitmq"
)

type (
	RMQConsumer interface {
		Consume(
			ctx context.Context,
			process func(ctx context.Context, msg rabbitmq.Delivery) error,
		) error
	}

	HandlerFunc func(ctx context.Context, id string, data any) error
)
