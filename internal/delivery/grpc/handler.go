package grpc_handler

import (
	"context"
	"github.com/ArturChopikian/grpc-server/configs"
	pb2 "github.com/ArturChopikian/grpc-server/internal/delivery/grpc/pb"
	"github.com/ArturChopikian/grpc-server/internal/models"
	"github.com/ArturChopikian/grpc-server/internal/usecase"
	"log"
)

type ProductsHandlerInterface interface {
	Fetch(ctx context.Context, req *pb2.FetchRequest) (*pb2.FetchResponse, error)
	List(ctx context.Context, req *pb2.ListRequest) (*pb2.ListResponse, error)
}

type productsHandler struct {
	productsUC usecase.ProductsUCInterface
	cfg        *configs.Config
	pb2.UnimplementedProductsServiceServer
}

func NewProductsHandler(useCase *usecase.UseCases, cfg *configs.Config) *productsHandler {
	return &productsHandler{
		productsUC: useCase.ProductsUC,
	}
}

func (s *productsHandler) Fetch(ctx context.Context, req *pb2.FetchRequest) (*pb2.FetchResponse, error) {

	err := s.productsUC.Fetch(ctx, req.GetUrl())
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb2.FetchResponse{
		Message: "Work",
	}, nil
}

func (s *productsHandler) List(ctx context.Context, req *pb2.ListRequest) (*pb2.ListResponse, error) {

	products, err := s.productsUC.List(ctx, req.GetOrderBy(), req.GetPageSize(), req.GetPageNumber())
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb2.ListResponse{
		Products:       models.ProductsToGrpc(products),
		NextPageNumber: req.GetPageNumber() + 1,
	}, nil
}
