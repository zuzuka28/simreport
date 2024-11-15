package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

type ProducerConfig struct {
	DSN          string `yaml:"dsn"`
	ExchangeName string `yaml:"exchangeName"`
	RoutingKey   string `yaml:"routingKey"`
	MaxRetries   int    `yaml:"maxRetries"`
}

type Producer struct {
	mu     sync.Mutex
	conn   *amqp091.Connection
	ch     *amqp091.Channel
	config ProducerConfig
}

func NewProducer(
	config ProducerConfig,
) (*Producer, error) {
	publisher := &Producer{
		mu:     sync.Mutex{},
		conn:   nil,
		ch:     nil,
		config: config,
	}

	if err := publisher.connect(); err != nil {
		return nil, err
	}

	if err := publisher.declareExchange(); err != nil {
		return nil, err
	}

	return publisher, nil
}

func (p *Producer) Publish(
	ctx context.Context,
	doc any,
) error {
	m, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	for {
		p.mu.Lock()
		err = p.ch.PublishWithContext(
			ctx,
			p.config.ExchangeName, // exhange
			p.config.RoutingKey,   // routing key
			false,                 // mandatory
			false,                 // immediate
			amqp091.Publishing{ //nolint:exhaustruct
				ContentType: "application/json",
				Body:        m,
			},
		)
		p.mu.Unlock()

		if err == nil {
			return nil
		}

		slog.Warn("error publishing message. attempting to reconnect...", "err", err)

		if err := p.reconnect(); err != nil {
			return fmt.Errorf("reconnect: %w", err)
		}
	}
}

func (p *Producer) connect() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var err error

	p.conn, err = amqp091.Dial(p.config.DSN)
	if err != nil {
		return fmt.Errorf("connecting to RabbitMQ: %w", err)
	}

	p.ch, err = p.conn.Channel()
	if err != nil {
		p.conn.Close()
		return fmt.Errorf("creating channel: %w", err)
	}

	return nil
}

func (p *Producer) close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.ch != nil {
		p.ch.Close()
	}

	if p.conn != nil {
		p.conn.Close()
	}
}

func (p *Producer) reconnect() error {
	p.close()

	var (
		err   error
		count int
	)

	for {
		if count > p.config.MaxRetries {
			return err
		}

		err = p.connect()
		if err == nil {
			slog.Info("successfully restored connection and channel.")
			return nil
		}

		time.Sleep(2 * time.Second) //nolint:mnd,gomnd

		count++
	}
}

func (p *Producer) declareExchange() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	err := p.ch.ExchangeDeclare(
		p.config.ExchangeName, // exhange
		"direct",              // kind
		true,                  // durable
		false,                 // delete when unused
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		return fmt.Errorf("declaring exchange: %w", err)
	}

	return nil
}
