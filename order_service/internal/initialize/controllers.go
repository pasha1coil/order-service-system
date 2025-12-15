package initialize

import (
	"order-service-system/order_service/internal/controllers/grpc/order_grpc_controller"

	"go.uber.org/zap"
)

type RpcControllersDeps struct {
	Logger   *zap.Logger
	Services *Services
}

type RpcControllers struct {
	OrderController *order_grpc_controller.OrderController
}

func NewRpcControllers(deps RpcControllersDeps) *RpcControllers {
	return &RpcControllers{
		OrderController: order_grpc_controller.NewOrderController(order_grpc_controller.Deps{
			Logger:       deps.Logger,
			OrderService: deps.Services.OrderServices,
		}),
	}
}
