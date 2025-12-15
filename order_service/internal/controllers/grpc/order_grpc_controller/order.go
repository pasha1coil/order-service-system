package order_grpc_controller

import (
	"context"
	"order-service-system/order_service/internal/service/order_service"
	orderpb "order-service-system/proto/order"

	"go.uber.org/zap"
)

type OrderController struct {
	orderpb.UnimplementedOrderServiceServer
	logger       *zap.Logger
	orderService *order_service.OrderService
}

type Deps struct {
	Logger       *zap.Logger
	OrderService *order_service.OrderService
}

func NewOrderController(deps Deps) *OrderController {
	if deps.Logger == nil {
		panic("deps.Logger is required")
	}

	if deps.OrderService == nil {
		panic("deps.Service is required")
	}

	return &OrderController{
		logger:       deps.Logger,
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
