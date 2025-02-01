package documentpipeline

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/document/internal/model"

	"github.com/nats-io/nats.go/jetstream"
)

const subject = "documentstatus"

type stageHandler struct {
	ss    StatusService
	stage Stage
	con   jetstream.Consumer
}

func newStageHandler(
	ctx context.Context,
	cm jetstream.ConsumerManager,
	s Stage,
	ss StatusService,
) (*stageHandler, error) {
	fullsubject := subject + "." + string(s.Trigger)

	con, err := cm.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{ //nolint:exhaustruct
		FilterSubject: fullsubject,
	})
	if err != nil {
		return nil, fmt.Errorf("create new consumer: %w", err)
	}

	return &stageHandler{
		ss:    ss,
		stage: s,
		con:   con,
	}, nil
}

func (h *stageHandler) Start(ctx context.Context) error {
	msgs, err := h.con.Messages()
	if err != nil {
		return fmt.Errorf("enter messages context: %w", err)
	}

	defer msgs.Stop()

	for {
		msg, err := msgs.Next()
		if err != nil {
			return fmt.Errorf("next message: %w", err)
		}

		if err := h.stage.Action.Serve(ctx, string(msg.Data())); err != nil {
			_ = msg.Nak()
			return fmt.Errorf("call action: %w", err)
		}

		if err := h.ss.Update(ctx, model.DocumentStatusUpdateCommand{
			ID:     string(msg.Data()),
			Status: h.stage.Next,
		}); err != nil {
			_ = msg.Nak()
			return fmt.Errorf("set next status: %w", err)
		}

		_ = msg.DoubleAck(ctx)
	}
}
