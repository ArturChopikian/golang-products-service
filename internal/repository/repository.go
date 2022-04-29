package repository

import (
	"context"
	"github.com/ArturChopikian/grpc-server/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductsReposInterface interface {
	Get(ctx context.Context, name string) (*models.Product, error)
	Create(ctx context.Context, product *models.Product) error
	UpdatePrice(ctx context.Context, id primitive.ObjectID, price float64) error
	List(ctx context.Context, orderBy map[string]int32, pageSize int32, pageNumber int32) ([]*models.Product, error)
}

type Repository struct {
	Products ProductsReposInterface
}

func NewRepository(coll *mongo.Collection) *Repository {
	return &Repository{
		Products: newProductsRepos(coll),
	}
}
