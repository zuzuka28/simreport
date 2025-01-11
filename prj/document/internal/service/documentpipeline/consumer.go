package documentpipeline

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	handlers []*stageHandler
}

func NewService(
	ctx context.Context,
	cm jetstream.ConsumerManager,
	ss StatusService,
	stages []Stage,
) (*Service, error) {
	handlers := make([]*stageHandler, 0, len(stages))

	for _, stage := range stages {
		handler, err := newStageHandler(ctx, cm, stage, ss)
		if err != nil {
			return nil, fmt.Errorf("new stage handler: %w", err)
		}

		handlers = append(handlers, handler)
	}

	return &Service{
		handlers: handlers,
	}, nil
}

func (s *Service) Start(ctx context.Context) error {
	eg, egCtx := errgroup.WithContext(ctx)

	for _, handler := range s.handlers {
		eg.Go(func() error {
			return handler.Start(egCtx)
		})
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("serve: %w", err)
	}

	return nil
}
