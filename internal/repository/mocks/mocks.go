package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"intresting_task/internal/models"
)

type ProductsReposMock struct {
	mock.Mock
}

func (p *ProductsReposMock) Get(ctx context.Context, name string) (*models.Product, error) {
	args := p.Called(ctx, name)

	return args.Get(0).(*models.Product), args.Error(1)
}

func (p *ProductsReposMock) Create(ctx context.Context, product *models.Product) error {
	args := p.Called(ctx, product)

	return args.Error(0)
}

func (p *ProductsReposMock) UpdatePrice(ctx context.Context, id primitive.ObjectID, price float64) error {
	args := p.Called(ctx, id, price)

	return args.Error(0)
}

func (p *ProductsReposMock) List(ctx context.Context, orderBy map[string]int32, pageSize int32, pageNumber int32) ([]*models.Product, error) {
	args := p.Called(ctx, orderBy, pageSize, pageNumber)

	return args.Get(0).([]*models.Product), args.Error(1)
}
