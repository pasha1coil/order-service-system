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
	return &Clients{
		OrderClient: clients.NewOrderClient(clients.OrderClientDeps{
			Logger:           deps.Logger,
			OrderServiceHost: deps.OrderServiceHost,
		}),
	}
}
