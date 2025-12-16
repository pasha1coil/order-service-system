package initialize

import (
	"order-service-system/order_service/internal/service/order_service"

	"go.uber.org/zap"
)

type Services struct {
	OrderServices *order_service.OrderService
}

type ServicesDeps struct {
	Logger       *zap.Logger
	Repositories *Repositories
	Clients      *Clients
}

func NewServices(deps ServicesDeps) *Services {
	if deps.Logger == nil {
		panic("logger must not be nil on <NewServices> of <initialize>")
	}
	if deps.Repositories == nil {
		panic("repositories must not be nil on <NewServices> of <initialize>")
	}
	if deps.Clients == nil {
		panic("clients must not be nil on <NewServices> of <initialize>")
	}
	return &Services{
		OrderServices: order_service.NewOrderService(order_service.Deps{
			Logger:     deps.Logger,
			OrderRepo:  deps.Repositories.OrderRepository,
			NatsClient: deps.Clients.NatsClient,
		}),
	}
}
