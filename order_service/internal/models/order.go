package models

import (
	orderpb "order-service-system/proto/order"
	"time"
)

type Order struct {
	OrderID     string      `bson:"order_id"`
	UserID      string      `bson:"user_id"`
	Items       []OrderItem `bson:"items"`
	TotalAmount float64     `bson:"total_amount"`
	Status      string      `bson:"status"`
	CreatedAt   time.Time   `bson:"created_at"`
	UpdatedAt   time.Time   `bson:"updated_at"`
}

type OrderItem struct {
	ProductID string  `bson:"product_id"`
	Quantity  int32   `bson:"quantity"`
	Price     float64 `bson:"price"`
}

type OrderCreatedEvent struct {
	OrderID     string
	UserID      string
	TotalAmount float64
	CreatedAt   time.Time
}

var AllowedStatuses = map[orderpb.OrderStatus]struct{}{
	orderpb.OrderStatus_PENDING:   {},
	orderpb.OrderStatus_PAID:      {},
	orderpb.OrderStatus_CANCELLED: {},
	orderpb.OrderStatus_FAILED:    {},
}
