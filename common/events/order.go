package events

type OrderCreatedPayload struct {
	OrderID     string  `json:"order_id"`
	UserID      string  `json:"user_id"`
	TotalAmount float64 `json:"total_amount"`
	CreatedAt   int64   `json:"created_at"`
}

type OrderPaidPayload struct {
	OrderID     string  `json:"order_id"`
	UserID      string  `json:"user_id"`
	TotalAmount float64 `json:"total_amount"`
	PaidAt      int64   `json:"paid_at"`
}

type OrderFailedPayload struct {
	OrderID  string `json:"order_id"`
	UserID   string `json:"user_id"`
	Reason   string `json:"reason"`
	FailedAt int64  `json:"failed_at"`
}
