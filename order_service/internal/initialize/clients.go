package initialize

import (
	"order-service-system/order_service/internal/clients/nats_client"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Clients struct {
	NatsClient *nats_client.Client
}

type ClientsDeps struct {
	Logger *zap.Logger
	Conn   *nats.Conn
}

func NewClients(deps ClientsDeps) *Clients {
	return &Clients{
		NatsClient: nats_client.NewClient(nats_client.Deps{
			Logger: deps.Logger,
			Conn:   deps.Conn,
		}),
	}
}
