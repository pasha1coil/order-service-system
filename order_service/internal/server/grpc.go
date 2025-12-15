package server

import (
	"context"
	"errors"
	"net"
	"order-service-system/order_service/internal/initialize"
	"order-service-system/proto/order"
	"time"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type DepsGRPC struct {
	Logger *zap.Logger
}

type GRPC struct {
	logger *zap.Logger
	grpc   *grpc.Server
}

func NewGRPC(deps DepsGRPC) (*GRPC, error) {
	if deps.Logger == nil {
		return nil, errors.New("logger is nil on <NewGRPC>")
	}

	grpcStreamInterceptor := grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(
		grpczap.StreamServerInterceptor(deps.Logger),
		grpcrecovery.StreamServerInterceptor(),
	))

	grpcUnaryInterceptor := grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
		grpczap.UnaryServerInterceptor(deps.Logger),
		grpcrecovery.UnaryServerInterceptor(),
	))

	return &GRPC{
		grpc:   grpc.NewServer(grpcStreamInterceptor, grpcUnaryInterceptor, grpc.ConnectionTimeout(5*time.Second)),
		logger: deps.Logger,
	}, nil
}

func (receiver *GRPC) Run(addr string) error {
	receiver.logger.Info("Starting GRPC Server", zap.String("host", addr))

	if err := receiver.listen(addr); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		receiver.logger.Error("GRPC Listen error", zap.Error(err))
		return err
	}
	return nil
}

func (receiver *GRPC) Stop(_ context.Context) error {
	receiver.grpc.GracefulStop()
	receiver.logger.Info("Shutting down GRPC server...")

	return nil
}

func (receiver *GRPC) Register(controllers *initialize.RpcControllers) *GRPC {
	order.RegisterOrderServiceServer(receiver.grpc, controllers.OrderController)
	// another...
	return receiver
}

func (receiver *GRPC) listen(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	return receiver.grpc.Serve(listener)
}
