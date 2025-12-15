package billing

import (
	"context"
	"encoding/json"
	"math/rand"
	"order-service-system/common/events"
	"order-service-system/proto/clients"
	"time"

	orderpb "order-service-system/proto/order"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

const (
	subjectOrderCreated = "order.created"
	queueBilling        = "billing-workers"
)

type Processor struct {
	logger      *zap.Logger
	natsConn    *nats.Conn
	orderClient *clients.OrderClient
	successRate float64
	rand        *rand.Rand
}

type Deps struct {
	Logger      *zap.Logger
	NatsConn    *nats.Conn
	OrderClient *clients.OrderClient
	SuccessRate float64
}

func NewProcessor(deps Deps) *Processor {
	if deps.Logger == nil {
		panic("logger is required")
	}
	if deps.NatsConn == nil {
		panic("nats connection is required")
	}
	if deps.OrderClient == nil {
		panic("order client is required")
	}

	successRate := deps.SuccessRate
	if successRate < 0 {
		successRate = 0
	}
	if successRate > 1 {
		successRate = 1
	}

	return &Processor{
		logger:      deps.Logger,
		natsConn:    deps.NatsConn,
		orderClient: deps.OrderClient,
		successRate: successRate,
		rand:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (receiver *Processor) Start(ctx context.Context) (*nats.Subscription, error) {
	sub, err := receiver.natsConn.QueueSubscribe(subjectOrderCreated, queueBilling, func(msg *nats.Msg) {
		receiver.handleMessage(ctx, msg)
	})
	if err != nil {
		return nil, err
	}
	if err := receiver.natsConn.Flush(); err != nil {
		return nil, err
	}

	receiver.logger.Info("listening for order.created", zap.String("subject", subjectOrderCreated), zap.String("queue", queueBilling))
	return sub, nil
}

func (receiver *Processor) handleMessage(ctx context.Context, msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		receiver.logger.Warn("context cancelled before processing")
		return
	default:
	}

	var payload events.OrderCreatedPayload
	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		receiver.logger.Error("failed to decode order.created", zap.Error(err))
		return
	}

	if payload.OrderID == "" || payload.UserID == "" {
		receiver.logger.Error("invalid payload", zap.Any("payload", payload))
		return
	}

	receiver.logger.Info("processing payment",
		zap.String("order_id", payload.OrderID),
		zap.String("user_id", payload.UserID),
		zap.Float64("amount", payload.TotalAmount))

	delay := 1000 + receiver.rand.Intn(1000) // 1000-2000ms
	time.Sleep(time.Duration(delay) * time.Millisecond)

	success := receiver.rand.Float64() <= receiver.successRate
	receiver.publishResult(ctx, payload, success)
}

func (receiver *Processor) publishResult(ctx context.Context, payload events.OrderCreatedPayload, success bool) {
	status := orderpb.OrderStatus_FAILED
	subject := "order.failed"
	var event any

	if success {
		status = orderpb.OrderStatus_PAID
		subject = "order.paid"
		event = events.OrderPaidPayload{
			OrderID:     payload.OrderID,
			UserID:      payload.UserID,
			TotalAmount: payload.TotalAmount,
			PaidAt:      time.Now().Unix(),
		}
	} else {
		event = events.OrderFailedPayload{
			OrderID:  payload.OrderID,
			UserID:   payload.UserID,
			Reason:   "payment declined",
			FailedAt: time.Now().Unix(),
		}
	}

	data, err := json.Marshal(event)
	if err != nil {
		receiver.logger.Error("failed to marshal event", zap.Error(err))
		return
	}

	if err := receiver.natsConn.Publish(subject, data); err != nil {
		receiver.logger.Error("failed to publish billing event",
			zap.String("subject", subject),
			zap.Error(err))
		receiver.logger.Warn("continuing despite publish error")
	} else {
		receiver.logger.Info("published event",
			zap.String("subject", subject),
			zap.String("order_id", payload.OrderID))
	}

	if _, err := receiver.orderClient.UpdateOrderStatus(ctx, payload.OrderID, status); err != nil {
		receiver.logger.Error("failed to update order status", zap.String("order_id", payload.OrderID), zap.Error(err))
		return
	}

	receiver.logger.Info("order status updated",
		zap.String("order_id", payload.OrderID),
		zap.String("status", status.String()))
}
