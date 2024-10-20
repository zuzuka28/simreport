package consumer

import (
	"context"
	"errors"
	"fmt"
	"simrep/internal/model"
	"simrep/pkg/rabbitmq"
)

var errUnprocessableType = errors.New("unprocessable notification type")

type Consumer struct {
	con      RMQConsumer
	handlers map[model.NotifyAction]HandlerFunc
}

func New(
	con RMQConsumer,
	handlers map[model.NotifyAction]HandlerFunc,
) *Consumer {
	return &Consumer{
		con:      con,
		handlers: handlers,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	cb := func(ctx context.Context, msg rabbitmq.Delivery) error {
		notif, err := parseDelivery(msg)
		if err != nil {
			return fmt.Errorf("parse notification: %w", err)
		}

		h, ok := c.handlers[model.NotifyAction(notif.Action)]
		if !ok {
			return fmt.Errorf("%w: %s", errUnprocessableType, notif.Action)
		}

		if err := h(ctx, notif.DocumentID, notif.UserData); err != nil {
			return fmt.Errorf("handle notify: %w", err)
		}

		return nil
	}

	if err := c.con.Consume(ctx, cb); err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	return nil
}
