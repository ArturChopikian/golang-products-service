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

// if page size equal zero, use default value instead of zero
const defaultPageSize = 100

// productRepos - define all methods for communicating with MongoDB collection
type productsRepos struct {
	conn *mongo.Collection
}

// Get - takes a name and return the product with this name
// or NotFoundProductError if product not found
// or error if some go wrong
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

// Create - take the product and insert it into collection
func (p *productsRepos) Create(ctx context.Context, product *models.Product) error {
	_, err := p.conn.InsertOne(ctx, &product)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("repos: Create: %v", err)
	}
	return nil
}

// UpdatePrice - take id and new price for product
// increase price_update by 1
// replace old price to new
// update time when price changed
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

// List - take orderBy map, pageSize, pageNumber
// ---
// orderBy map represent all fields for ordering look like:
// "order_by": {
//    "price": 1
//  }
// Where "price" is name of field and "1" it is determines ascending(1)/descending(-1) sort
// ---
// pageSize is maximum products per one page
// ---
// pageNumber is current page and represented how many pages need to skip
// ---
// Return list of product's pointers
// or error if something went wrong
func (p *productsRepos) List(ctx context.Context, orderBy map[string]int32, pageSize int32, pageNumber int32) ([]*models.Product, error) {
	opts := &options.FindOptions{}

	// if page number == 0 we don't need to pass some products
	if pageNumber == 0 {
		opts.SetSkip(0)
	} else {
		opts.SetSkip(int64(pageSize) * int64(pageNumber))
	}

	// set up all order settings
	for key, value := range orderBy {
		opts.SetSort(bson.D{{key, value}})
	}

	// if page size == 0 we have default value for it
	if pageSize == 0 {
		opts.SetLimit(defaultPageSize)
	} else {
		opts.SetLimit(int64(pageSize))
	}

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
		result = append(result, p)
	}

	if err := cur.Err(); err != nil {
		log.Println("repos: List: error while get Err from cursor:", err)
		return nil, err
	}
	return result, nil
}

// newProductsRepos - return new productsRepos
func newProductsRepos(conn *mongo.Collection) *productsRepos {
	return &productsRepos{
		conn: conn,
	}
}
