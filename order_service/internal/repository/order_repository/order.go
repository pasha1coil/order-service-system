package order_repository

import (
	"context"
	"errors"
	"order-service-system/order_service/internal/models"
	"order-service-system/order_service/internal/pj_errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type OrderRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

type Deps struct {
	Collection *mongo.Collection
}

func NewOrderRepository(ctx context.Context, deps Deps) (*OrderRepository, error) {
	if deps.Collection == nil {
		panic("deps.Collection cannot be nil")
	}

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "order_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	if _, err := deps.Collection.Indexes().CreateOne(ctx, indexModel); err != nil {
		return nil, err
	}

	return &OrderRepository{
		collection: deps.Collection,
	}, nil
}

func (receiver *OrderRepository) Create(ctx context.Context, order models.Order) error {
	_, err := receiver.collection.InsertOne(ctx, order)
	return err
}

func (receiver *OrderRepository) Get(ctx context.Context, orderID string) (models.Order, error) {
	var doc models.Order
	err := receiver.collection.FindOne(ctx, bson.M{"order_id": orderID}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return models.Order{}, pj_errors.ErrNotFound
	}
	return doc, err
}

func (receiver *OrderRepository) UpdateStatus(ctx context.Context, orderID string, status string) (models.Order, error) {
	var doc models.Order
	res := receiver.collection.FindOneAndUpdate(ctx,
		bson.M{"order_id": orderID},
		bson.M{"$set": bson.M{"status": status, "updated_at": time.Now().UTC()}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Order{}, pj_errors.ErrNotFound
		}
		return models.Order{}, err
	}
	if err := res.Decode(&doc); err != nil {
		return models.Order{}, err
	}
	return doc, nil
}
