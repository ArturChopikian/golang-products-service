package models

import (
	"github.com/ArturChopikian/grpc-server/internal/delivery/grpc/pb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Product struct {
	Id           primitive.ObjectID `bson:"_id"`
	Name         string             `bson:"name"`
	Price        float64            `bson:"price"`
	Updated      time.Time          `bson:"updated"`
	PriceUpdates uint32             `bson:"price_updates"`
}

func ProductToGrpc(p *Product) *pb.Product {
	return &pb.Product{
		Id:           p.Id.Hex(),
		Name:         p.Name,
		Price:        p.Price,
		Updated:      timeToTimestamp(p.Updated),
		PriceUpdates: p.PriceUpdates,
	}
}

func ProductsToGrpc(p []*Product) []*pb.Product {
	var result []*pb.Product

	for _, product := range p {
		result = append(result, ProductToGrpc(product))
	}

	return result
}

func timeToTimestamp(t time.Time) *timestamppb.Timestamp {
	return &timestamppb.Timestamp{
		Seconds: int64(t.Second()),
		Nanos:   int32(t.Nanosecond()),
	}
}
