package producer

import (
	"context"
	"fmt"
	"simrep/internal/model"
)

type Producer struct {
	pub Publisher
}

func New(
	rmq Publisher,
) *Producer {
	return &Producer{
		pub: rmq,
	}
}

func (p *Producer) Notify(
	ctx context.Context,
	documentID string,
	action model.NotifyAction,
	userdata any,
) error {
	notif := notification{
		DocumentID: documentID,
		Action:     string(action),
		UserData:   userdata,
	}

	if err := p.pub.Publish(ctx, notif); err != nil {
		return fmt.Errorf("can't publish notification: %w", err)
	}

	return nil
}
