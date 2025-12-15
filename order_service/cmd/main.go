package main

import (
	"context"
	"fmt"
	"order-service-system/order_service/internal/app"
	"order-service-system/order_service/internal/initialize"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	config, err := initialize.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err = app.Run(ctx, *config, logger); err != nil {
		logger.Fatal("App exited with error", zap.Error(err))
	}
}
