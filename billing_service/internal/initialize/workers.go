package initialize

import (
	"order-service-system/billing_service/internal/workers/billing"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Workers struct {
	BillingProcessor *billing.Processor
}

type WorkersDeps struct {
	Logger      *zap.Logger
	Clients     *Clients
	NatsConn    *nats.Conn
	SuccessRate float64
}

func NewWorkers(deps WorkersDeps) *Workers {
	if deps.Logger == nil {
		panic("logger must not be nil on <NewWorkers> of <initialize>")
	}
	if deps.NatsConn == nil {
		panic("nats connection must not be nil on <NewWorkers> of <initialize>")
	}
	if deps.Clients == nil {
		panic("clients must not be nil on <NewWorkers> of <initialize>")
	}
	return &Workers{
		BillingProcessor: billing.NewProcessor(billing.Deps{
			Logger:      deps.Logger,
			NatsConn:    deps.NatsConn,
			OrderClient: deps.Clients.OrderClient,
			SuccessRate: deps.SuccessRate,
		}),
	}
}
