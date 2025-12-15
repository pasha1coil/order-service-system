package clients

import (
	"context"
	"order-service-system/proto/order"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderClientDeps struct {
	Logger           *zap.Logger
	OrderServiceHost string
}

type OrderClient struct {
	logger           *zap.Logger
	orderServiceHost string
}

func NewOrderClient(deps OrderClientDeps) *OrderClient {
	if deps.Logger == nil {
		panic("logger must not be nil")
	}
	if deps.OrderServiceHost == "" {
		panic("orderServiceHost must not be empty")
	}
	return &OrderClient{
		logger:           deps.Logger,
		orderServiceHost: deps.OrderServiceHost,
	}
}

func (receiver *OrderClient) UpdateOrderStatus(ctx context.Context, orderID string, status order.OrderStatus) (*order.UpdateOrderStatusResponse, error) {
	connection, err := grpc.Dial(receiver.orderServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		receiver.logger.Error("failed to connect on <UpdateOrderStatus> of <OrderClient>", zap.Error(err), zap.String("order host", receiver.orderServiceHost))
		return nil, err
	}
	defer func() {
		if closeErr := connection.Close(); closeErr != nil {
			receiver.logger.Error("failed to close connection on <UpdateOrderStatus> of <OrderClient>", zap.Error(closeErr))
		}
	}()

	client := order.NewOrderServiceClient(connection)

	response, err := client.UpdateOrderStatus(ctx, &order.UpdateOrderStatusRequest{
		OrderId: orderID,
		Status:  status,
	})
	if err != nil {
		receiver.logger.Error("failed Update Order Status on <UpdateOrderStatus> of <OrderClient>", zap.Error(err))
		return nil, err
	}

	return response, nil
}
