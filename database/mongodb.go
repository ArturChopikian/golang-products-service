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
		SetServerAPIOptions(serverAPIOptions).
		SetMaxPoolSize(1024)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	return client.Database("products-service").Collection("products"), nil
}
