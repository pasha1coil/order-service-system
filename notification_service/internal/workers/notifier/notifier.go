package notifier

import (
	"context"
	"encoding/json"
	"order-service-system/common/events"
	"order-service-system/proto/clients"

	orderpb "order-service-system/proto/order"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

const (
	subjectOrderPaid   = "order.paid"
	subjectOrderFailed = "order.failed"
	queueNotification  = "notification-workers"
)

type Notifier struct {
	logger      *zap.Logger
	natsConn    *nats.Conn
	orderClient *clients.OrderClient
}

type Deps struct {
	Logger      *zap.Logger
	NatsConn    *nats.Conn
	OrderClient *clients.OrderClient
}

func New(deps Deps) *Notifier {
	if deps.Logger == nil {
		panic("logger is required")
	}
	if deps.NatsConn == nil {
		panic("nats connection is required")
	}
	if deps.OrderClient == nil {
		panic("order client is required")
	}
	return &Notifier{
		logger:      deps.Logger,
		natsConn:    deps.NatsConn,
		orderClient: deps.OrderClient,
	}
}

func (receiver *Notifier) Start(ctx context.Context) ([]*nats.Subscription, error) {
	subPaid, err := receiver.natsConn.QueueSubscribe(subjectOrderPaid, queueNotification, func(msg *nats.Msg) {
		receiver.handlePaid(ctx, msg)
	})
	if err != nil {
		return nil, err
	}

	subFailed, err := receiver.natsConn.QueueSubscribe(subjectOrderFailed, queueNotification, func(msg *nats.Msg) {
		receiver.handleFailed(ctx, msg)
	})
	if err != nil {
		return nil, err
	}

	if err := receiver.natsConn.Flush(); err != nil {
		return nil, err
	}

	receiver.logger.Info("listening for payment events", zap.String("paid", subjectOrderPaid), zap.String("failed", subjectOrderFailed), zap.String("queue", queueNotification))
	return []*nats.Subscription{subPaid, subFailed}, nil
}

func (receiver *Notifier) handlePaid(ctx context.Context, msg *nats.Msg) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	var payload events.OrderPaidPayload
	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		receiver.logger.Error("failed to decode paid payload", zap.Error(err))
		return
	}

	if _, err := receiver.orderClient.UpdateOrderStatus(ctx, payload.OrderID, orderpb.OrderStatus_PAID); err != nil {
		receiver.logger.Error("failed to update order status", zap.String("order_id", payload.OrderID), zap.Error(err))
		return
	}

	receiver.logger.Info("notified user about payment", zap.String("order_id", payload.OrderID), zap.String("user_id", payload.UserID))
}

func (receiver *Notifier) handleFailed(ctx context.Context, msg *nats.Msg) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	var payload events.OrderFailedPayload
	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		receiver.logger.Error("failed to decode failed payload", zap.Error(err))
		return
	}

	if _, err := receiver.orderClient.UpdateOrderStatus(ctx, payload.OrderID, orderpb.OrderStatus_FAILED); err != nil {
		receiver.logger.Error("failed to update order status", zap.String("order_id", payload.OrderID), zap.Error(err))
		return
	}

	receiver.logger.Info("notified user about failure", zap.String("order_id", payload.OrderID), zap.String("user_id", payload.UserID), zap.String("reason", payload.Reason))
}
