package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"simrep/internal/model"
	"sync"
	"time"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	DSN        string `yaml:"dsn"`
	QueueName  string `yaml:"queueName"`
	MaxRetries int    `yaml:"maxRetries"`
}

type Consumer struct {
	conn   *amqp091.Connection
	ch     *amqp091.Channel
	mu     sync.Mutex
	config Config

	ds DocumentService
	as AnalyzeService
}

func New(
	config Config,
	dr DocumentService,
	as AnalyzeService,
) (*Consumer, error) {
	c := &Consumer{
		conn:   nil,
		ch:     nil,
		mu:     sync.Mutex{},
		config: config,
		ds:     dr,
		as:     as,
	}

	if err := c.connect(); err != nil {
		return nil, err
	}

	if err := c.declareQueue(); err != nil {
		return nil, fmt.Errorf("declare queue: %w", err)
	}

	return c, nil
}

func (c *Consumer) Start(ctx context.Context) error {
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
			var message analyzeRequestMessage

			if err := json.Unmarshal(msg.Body, &message); err != nil {
				return fmt.Errorf("unmarshal message: %w", err)
			}

			if err := c.processMessage(ctx, message); err != nil {
				return fmt.Errorf("process message: %w", err)
			}

			if err := msg.Ack(false); err != nil {
				return fmt.Errorf("ack message: %w", err)
			}
		}
	}
}

func (c *Consumer) processMessage(
	ctx context.Context,
	msg analyzeRequestMessage,
) error {
	slog.Info("got message", "msg", msg)

	file, err := c.ds.FetchParsedFile(ctx, model.ParsedDocumentFileQuery{
		DocumentID: msg.DocumentID,
	})
	if err != nil {
		return fmt.Errorf("fetch file: %w", err)
	}

	analyzed, err := c.as.Analyze(ctx, file)
	if err != nil {
		return fmt.Errorf("analyze file: %w", err)
	}

	if err := c.as.Save(ctx, model.AnalyzedDocumentSaveCommand{
		Item: analyzed,
	}); err != nil {
		return fmt.Errorf("save document: %w", err)
	}

	slog.Info("message processed", "msg", msg)

	return nil
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
