package consumer

import (
	"encoding/json"
	"fmt"
	"simrep/pkg/rabbitmq"
)

func parseDelivery(msg rabbitmq.Delivery) (*notification, error) {
	var notif notification
	if err := json.Unmarshal(msg.Body, &notif); err != nil {
		return nil, fmt.Errorf("unmarshal notification: %w", err)
	}

	return &notif, nil
}
