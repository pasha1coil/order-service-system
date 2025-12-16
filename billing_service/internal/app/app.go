package app

import (
	"context"
	"errors"
	"fmt"
	"order-service-system/billing_service/internal/initialize"
	"order-service-system/common/closer"
	"order-service-system/common/nats"
	"time"

	"go.uber.org/zap"
)

func Run(ctx context.Context, config initialize.Config, logger *zap.Logger) error {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("recovered from panic on <Run> of <app>", zap.Any("error", r))
		}
	}()

	logger.Info("service start on <Run> of <app>", zap.String("service", "billing"), zap.Any("config", config))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	shutdownGroup := closer.NewCloserGroup()

	natsConn, err := nats.Connect(config.ExternalCfg.NatsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect nats: %w", err)
	}

	clients := initialize.NewClients(initialize.ClientsDeps{
		Logger:           logger,
		OrderServiceHost: config.OrderServiceHost,
	})

	workers := initialize.NewWorkers(initialize.WorkersDeps{
		Logger:      logger,
		Clients:     clients,
		SuccessRate: config.PaymentSuccessRate,
		NatsConn:    natsConn,
	})

	subscription, err := workers.BillingProcessor.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed to subscribe to order.created: %w", err)
	}

	shutdownGroup.Add(closer.CloserFunc(func(ctx context.Context) error {
		return subscription.Drain()
	}))
	shutdownGroup.Add(closer.CloserFunc(func(ctx context.Context) error {
		return natsConn.Drain()
	}))

	<-ctx.Done()

	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer timeoutCancel()
	if err := shutdownGroup.Call(timeoutCtx); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Error("shutdown timed out on <Run> of <app>", zap.Error(err))
		} else {
			logger.Error("failed to shutdown services gracefully on <Run> of <app>", zap.Error(err))
		}
		return err
	}

	logger.Info("service stopped on <Run> of <app>", zap.String("service", "billing"))
	return nil
}
