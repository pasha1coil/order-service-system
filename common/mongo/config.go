package mongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Configuration struct {
	URL          string `env:"MONGO_URL"`
	DatabaseName string `env:"MONGO_DB_NAME"`
}

type RequestSettings struct {
	Driver  *mongo.Collection
	Options *options.FindOptions
	Filter  primitive.M
}
