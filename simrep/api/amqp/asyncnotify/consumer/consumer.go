package consumer

import (
	"context"
	"fmt"
	"simrep/pkg/rabbitmq"
)

type Consumer struct {
	con     RMQConsumer
	handler Handler
}

func New(
	con RMQConsumer,
	handler Handler,
) *Consumer {
	return &Consumer{
		con:     con,
		handler: handler,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	cb := func(ctx context.Context, msg rabbitmq.Delivery) error {
		notif, err := parseDelivery(msg)
		if err != nil {
			return fmt.Errorf("parse notification: %w", err)
		}

		if err := c.handler.Serve(ctx, notif.DocumentID, notif.UserData); err != nil {
			return fmt.Errorf("handle notify: %w", err)
		}

		return nil
	}

	if err := c.con.Consume(ctx, cb); err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	return nil
}
