package initialize

import (
	"log"
	"order-service-system/common/mongo"
	"order-service-system/common/nats"

	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
)

type Config struct {
	GrpcURL     string `env:"GRPC_URL"`
	ExternalCfg ExternalCfg
}

type ExternalCfg struct {
	MongoConfig mongo.Configuration
	NatsConfig  nats.Configuration
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
