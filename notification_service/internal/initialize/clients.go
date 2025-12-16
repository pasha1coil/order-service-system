package initialize

import (
	"order-service-system/proto/clients"

	"go.uber.org/zap"
)

type Clients struct {
	OrderClient *clients.OrderClient
}

type ClientsDeps struct {
	Logger           *zap.Logger
	OrderServiceHost string
}

func NewClients(deps ClientsDeps) *Clients {
	if deps.Logger == nil {
		panic("logger must not be nil on <NewClients> of <initialize>")
	}
	if deps.OrderServiceHost == "" {
		panic("order service host must not be empty on <NewClients> of <initialize>")
	}
	return &Clients{
		OrderClient: clients.NewOrderClient(clients.OrderClientDeps{
			Logger:           deps.Logger,
			OrderServiceHost: deps.OrderServiceHost,
		}),
	}
}
