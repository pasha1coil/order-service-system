package initialize

import (
	"order-service-system/notification_service/internal/workers/notifier"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Workers struct {
	Notifier *notifier.Notifier
}

type WorkersDeps struct {
	Logger   *zap.Logger
	NatsConn *nats.Conn
	Clients  *Clients
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
		Notifier: notifier.New(notifier.Deps{
			Logger:      deps.Logger,
			NatsConn:    deps.NatsConn,
			OrderClient: deps.Clients.OrderClient,
		}),
	}
}
