package server

import (
	"fmt"
	"github.com/ArturChopikian/grpc-server/configs"
	grpc_handler "github.com/ArturChopikian/grpc-server/internal/delivery/grpc"
	"github.com/ArturChopikian/grpc-server/internal/delivery/grpc/pb"
	"github.com/ArturChopikian/grpc-server/internal/repository"
	"github.com/ArturChopikian/grpc-server/internal/usecase"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type ProductServers struct {
	handler   *grpc_handler.ProductsHandler
	cfg       *configs.Config
	mongoColl *mongo.Collection
	server    *grpc.Server
	lis       net.Listener
}

func NewProductsServer(cfg *configs.Config, coll *mongo.Collection) (*ProductServers, error) {

	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	lis, err := net.Listen(cfg.Server.Network, address)
	if err != nil {
		return &ProductServers{}, err
	}

	return &ProductServers{
		cfg:       cfg,
		mongoColl: coll,
		server:    grpc.NewServer(),
		lis:       lis,
	}, nil
}

func (ps *ProductServers) MapHandler() {
	repositories := repository.NewRepository(ps.mongoColl)
	useCases := usecase.NewUseCases(repositories)
	handlers := grpc_handler.NewProductsHandler(useCases, ps.cfg)

	pb.RegisterProductsServiceServer(ps.server, handlers)
	reflection.Register(ps.server)
}

func (ps *ProductServers) Run() error {
	return ps.server.Serve(ps.lis)
}

func (ps *ProductServers) Stop() {

	if err := ps.lis.Close(); err != nil {
		log.Printf("ERROR: %v\n", err)
	}
	ps.server.Stop()
}
