package mongo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConnectDeps struct {
	Configuration *Configuration
	Timeout       time.Duration
}

func Connect(ctx context.Context, deps *ConnectDeps) (*mongo.Database, error) {
	if deps == nil {
		return nil, errors.New("arguments are empty")
	}

	connectionOptions := options.Client().
		ApplyURI(deps.Configuration.URL)

	ticker := time.NewTicker(1 * time.Second)
	timeoutExceeded := time.After(deps.Timeout)

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			connection, err := mongo.Connect(ctx, connectionOptions)

			if err == nil {
				err = connection.Ping(ctx, nil)
				if err == nil {
					return connection.Database(deps.Configuration.DatabaseName), nil
				}
				log.Printf("failed to ping the database <%s>: %s", deps.Configuration.URL, err.Error())
			}

			log.Printf("failed to connect to db <%s>: %s", deps.Configuration.URL, err.Error())
		case <-timeoutExceeded:
			return nil, fmt.Errorf("db connection <%s> failed after %d timeout", deps.Configuration.URL, deps.Timeout)
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
