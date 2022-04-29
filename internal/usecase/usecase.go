package usecase

import (
	"context"
	"github.com/ArturChopikian/grpc-server/internal/models"
	"github.com/ArturChopikian/grpc-server/internal/repository"
)

type ProductsUCInterface interface {
	Fetch(ctx context.Context, url string) error
	List(ctx context.Context, orderBy map[string]int32, pageSize int32, pageNumber int32) ([]*models.Product, error)
}

type UseCases struct {
	ProductsUC ProductsUCInterface
}

func NewUseCases(repos *repository.Repository) *UseCases {
	return &UseCases{ProductsUC: newProductUC(repos.Products)}
}
