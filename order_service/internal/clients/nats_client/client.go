package nats_client

import (
	"encoding/json"
	"fmt"
	"order-service-system/common/events"
	"order-service-system/order_service/internal/models"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Client struct {
	conn   *nats.Conn
	logger *zap.Logger
}

type Deps struct {
	Logger *zap.Logger
	Conn   *nats.Conn
}

func NewClient(deps Deps) *Client {
	if deps.Conn == nil {
		panic("nats connection required")
	}
	return &Client{
		logger: deps.Logger,
		conn:   deps.Conn,
	}
}

func (receiver *Client) PublishOrderCreated(event models.OrderCreatedEvent) error {
	payload := events.OrderCreatedPayload{
		OrderID:     event.OrderID,
		UserID:      event.UserID,
		TotalAmount: event.TotalAmount,
		CreatedAt:   event.CreatedAt.Unix(),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	if err := receiver.conn.Publish("order.created", data); err != nil {
		return fmt.Errorf("publish event: %w", err)
	}

	receiver.logger.Info("published event",
		zap.String("subject", "order.created"),
		zap.String("order_id", event.OrderID),
		zap.Time("created_at", event.CreatedAt),
	)
	return nil
}
