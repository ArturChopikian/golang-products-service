package grpc_handler

import (
	"context"
	"github.com/ArturChopikian/grpc-server/configs"
	"github.com/ArturChopikian/grpc-server/internal/delivery/grpc/pb"
	"github.com/ArturChopikian/grpc-server/internal/models"
	"github.com/ArturChopikian/grpc-server/internal/usecase"
	"log"
)

// ProductsHandlerInterface - represent productsHandler logic
type ProductsHandlerInterface interface {
	Fetch(ctx context.Context, req *pb.FetchRequest) (*pb.FetchResponse, error)
	List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error)
}

// productsHandler - implement handlers for ProductsService
type productsHandler struct {
	productsUC usecase.ProductsUCInterface
	cfg        *configs.Config
	pb.UnimplementedProductsServiceServer
}

// NewProductsHandler - return pointer of productsHandler
func NewProductsHandler(useCase *usecase.UseCases, cfg *configs.Config) *productsHandler {
	return &productsHandler{
		productsUC: useCase.ProductsUC,
	}
}

// Fetch - get data from external URL and after create new or update exists product
// take request with external URL and transmit to the next layer (useCases)
// return - pb.FetchResponse with message "work" or error
func (s *productsHandler) Fetch(ctx context.Context, req *pb.FetchRequest) (*pb.FetchResponse, error) {

	err := s.productsUC.Fetch(ctx, req.GetUrl())
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.FetchResponse{
		Message: "Work",
	}, nil
}

// List - implement endless scroll
// take pb.ListRequest with paging and ordering params
// return list of products and number of the next page or error
func (s *productsHandler) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {

	products, err := s.productsUC.List(ctx, req.GetOrderBy(), req.GetPageSize(), req.GetPageNumber())
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ListResponse{
		Products:       models.ProductsToGrpc(products),
		NextPageNumber: req.GetPageNumber() + 1,
	}, nil
}
