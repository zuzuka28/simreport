package rabbitmq

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

type ConsumerConfig struct {
	DSN        string `yaml:"dsn"`
	QueueName  string `yaml:"queueName"`
	MaxRetries int    `yaml:"maxRetries"`
}

type Consumer struct {
	conn   *amqp091.Connection
	ch     *amqp091.Channel
	mu     sync.Mutex
	config ConsumerConfig
}

func NewConsumer(
	config ConsumerConfig,
) (*Consumer, error) {
	c := &Consumer{
		conn:   nil,
		ch:     nil,
		mu:     sync.Mutex{},
		config: config,
	}

	if err := c.connect(); err != nil {
		return nil, err
	}

	if err := c.declareQueue(); err != nil {
		return nil, fmt.Errorf("declare queue: %w", err)
	}

	return c, nil
}

func (c *Consumer) Consume(
	ctx context.Context,
	process func(ctx context.Context, msg Delivery) error,
) error {
	slog.Info("Waiting for messages", "queue", c.config.QueueName)

	for {
		msgs, err := c.ch.ConsumeWithContext(
			ctx,
			c.config.QueueName, // queue
			"",                 // consumer
			false,              // auto-ack
			false,              // exclusive
			false,              // no-local
			false,              // no-wait
			nil,                // args
		)
		if err != nil {
			if err := c.reconnect(); err != nil {
				return err
			}
		}

		for msg := range msgs {
			if err := process(ctx, Delivery(msg)); err != nil {
				return fmt.Errorf("process message: %w", err)
			}

			if err := msg.Ack(false); err != nil {
				return fmt.Errorf("ack message: %w", err)
			}
		}
	}
}

func (c *Consumer) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var err error

	c.conn, err = amqp091.Dial(c.config.DSN)
	if err != nil {
		return fmt.Errorf("connect to RabbitMQ: %w", err)
	}

	c.ch, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("open a channel: %w", err)
	}

	return nil
}

func (c *Consumer) close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ch != nil {
		c.ch.Close()
	}

	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Consumer) reconnect() error {
	c.close()

	var (
		err   error
		count int
	)

	for {
		if count > c.config.MaxRetries {
			return err
		}

		err = c.connect()
		if err == nil {
			slog.Info("successfully restored connection and channel.")
			return nil
		}

		time.Sleep(2 * time.Second) //nolint:mnd,gomnd

		count++
	}
}

func (c *Consumer) declareQueue() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, err := c.ch.QueueDeclare(
		c.config.QueueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("declaring queue: %w", err)
	}

	return nil
}
