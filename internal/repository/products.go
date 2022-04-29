package repository

import (
	"context"
	"fmt"
	"github.com/ArturChopikian/grpc-server/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type productsRepos struct {
	conn *mongo.Collection
}

func (p *productsRepos) Get(ctx context.Context, name string) (*models.Product, error) {
	filter := bson.M{"name": name}
	product := &models.Product{}

	result := p.conn.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.NotFoundProductError
		}
		log.Println(err)
		return nil, fmt.Errorf("repos: Get: finding product by name: %v", err)
	}

	if err := result.Decode(&product); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("repos: Get: decoding finded product: %v", err)
	}

	return product, nil
}

func (p *productsRepos) Create(ctx context.Context, product *models.Product) error {
	_, err := p.conn.InsertOne(ctx, &product)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("repos: Create: %v", err)
	}
	return nil
}

func (p *productsRepos) UpdatePrice(ctx context.Context, id primitive.ObjectID, price float64) error {
	filter := bson.M{"_id": id}
	// update price, updated time and increase update price counter by 1
	update := bson.D{
		{"$inc", bson.D{{"price_updates", 1}}},
		{"$set", bson.D{{"price", price}}},
		{"$set", bson.D{{"updated", time.Now()}}},
	}

	res, err := p.conn.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("repos: IncreasePriceUpdates: %v", err)
	}
	if res.ModifiedCount == 0 {
		log.Println(err)
		return fmt.Errorf("repos: IncreasePriceUpdates: nothing updated")
	}
	return nil
}

func (p *productsRepos) List(ctx context.Context, orderBy map[string]int32, pageSize int32, pageNumber int32) ([]*models.Product, error) {
	opts := &options.FindOptions{}

	if pageNumber == 0 {
		opts.SetSkip(0)
	} else {
		opts.SetSkip(int64(pageSize) * int64(pageNumber))
	}

	for key, value := range orderBy {
		opts.SetSort(bson.D{{key, value}})
	}

	opts.SetLimit(int64(pageSize))

	cur, err := p.conn.Find(ctx, bson.D{}, opts)
	if err != nil {
		log.Println("repos: List: error while finding:", err)
		return nil, err
	}

	var result []*models.Product

	for cur.Next(ctx) {
		p := &models.Product{}

		err := cur.Decode(p)
		if err != nil {
			log.Println("repos: List: error while decoding:", err)
			return nil, err
		}

		fmt.Println(p)
		result = append(result, p)
	}

	fmt.Println("len result:", len(result))
	fmt.Println("cap result:", cap(result))
	fmt.Println(result)

	if err := cur.Err(); err != nil {
		log.Println("repos: List: error while get Err from cursor:", err)
		return nil, err
	}
	return result, nil
}

func newProductsRepos(conn *mongo.Collection) *productsRepos {
	return &productsRepos{
		conn: conn,
	}
}
