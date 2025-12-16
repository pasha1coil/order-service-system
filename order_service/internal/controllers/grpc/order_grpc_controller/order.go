package order_grpc_controller

import (
	"context"
	"order-service-system/order_service/internal/service/order_service"
	orderpb "order-service-system/proto/order"
)

type OrderController struct {
	orderpb.UnimplementedOrderServiceServer
	orderService *order_service.OrderService
}

type Deps struct {
	OrderService *order_service.OrderService
}

func NewOrderController(deps Deps) *OrderController {
	if deps.OrderService == nil {
		panic("deps.Service is required")
	}

	return &OrderController{
		orderService: deps.OrderService,
	}
}

func (receiver *OrderController) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	order, err := receiver.orderService.CreateOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	return &orderpb.CreateOrderResponse{Order: order}, nil
}

func (receiver *OrderController) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.GetOrderResponse, error) {
	order, err := receiver.orderService.GetOrder(ctx, req.GetOrderId())
	if err != nil {
		return nil, err
	}
	return &orderpb.GetOrderResponse{Order: order}, nil
}

func (receiver *OrderController) UpdateOrderStatus(ctx context.Context, req *orderpb.UpdateOrderStatusRequest) (*orderpb.UpdateOrderStatusResponse, error) {
	order, err := receiver.orderService.UpdateOrderStatus(ctx, req.GetOrderId(), req.GetStatus())
	if err != nil {
		return nil, err
	}
	return &orderpb.UpdateOrderStatusResponse{Order: order}, nil
}
