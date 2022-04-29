package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewMongoDBCollection() (*mongo.Collection, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)

	clientOptions := options.Client().
		ApplyURI("mongodb+srv://arturchopikian:cCgBuL7CWzz5htdQ@cluster0.7b913.mongodb.net/test?retryWrites=true&w=majority").
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
