package initialize

import (
	"log"
	"order-service-system/common/nats"

	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
)

type Config struct {
	PaymentSuccessRate float64 `env:"PAYMENT_SUCCESS_RATE"`
	OrderServiceHost   string  `env:"ORDER_SERVICE_HOST"`
	ExternalCfg        ExternalCfg
}

type ExternalCfg struct {
	NatsConfig nats.Configuration
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	var config Config
	if err := env.Parse(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
