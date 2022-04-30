package database

import (
	"context"
	"github.com/ArturChopikian/grpc-server/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewMongoDBCollection(cfg *configs.Config) (*mongo.Collection, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)

	clientOptions := options.Client().
		ApplyURI(cfg.MongoDB.URI).
		SetServerAPIOptions(serverAPIOptions)

	// timeout in seconds
	timeout := time.Duration(cfg.MongoDB.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	return client.Database(cfg.MongoDB.Database).Collection(cfg.MongoDB.Collection), nil
}
