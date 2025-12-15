package initialize

import (
	"context"
	"order-service-system/order_service/internal/repository/order_repository"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	OrderRepository *order_repository.OrderRepository
}

type RepositoriesDeps struct {
	MongoDB *mongo.Database
}

func NewRepositories(ctx context.Context, deps RepositoriesDeps) (*Repositories, error) {
	orderRepo, err := order_repository.NewOrderRepository(ctx, order_repository.Deps{
		Collection: deps.MongoDB.Collection("order"),
	})
	if err != nil {
		return nil, err
	}
	return &Repositories{
		OrderRepository: orderRepo,
	}, nil
}
