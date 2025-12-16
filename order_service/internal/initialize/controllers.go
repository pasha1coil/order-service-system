package initialize

import (
	"order-service-system/order_service/internal/controllers/grpc/order_grpc_controller"
)

type RpcControllersDeps struct {
	Services *Services
}

type RpcControllers struct {
	OrderController *order_grpc_controller.OrderController
}

func NewRpcControllers(deps RpcControllersDeps) *RpcControllers {
	return &RpcControllers{
		OrderController: order_grpc_controller.NewOrderController(order_grpc_controller.Deps{
			OrderService: deps.Services.OrderServices,
		}),
	}
}
