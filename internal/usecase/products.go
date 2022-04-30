package usecase

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/ArturChopikian/grpc-server/internal/models"
	"github.com/ArturChopikian/grpc-server/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// productUC - define business logic for products handlers
type productUC struct {
	productsRepos repository.ProductsReposInterface
}

// List - take orderBy, pageSize, pageNumber
// and call List method from repository
// return list of product's pointers or error
func (uc *productUC) List(ctx context.Context, orderBy map[string]int32, pageSize int32, pageNumber int32) ([]*models.Product, error) {

	return uc.productsRepos.List(ctx, orderBy, pageSize, pageNumber)
}

// Fetch - take URL with external csv file
// we have the pipeline
//
//			->check->
//
//			->check->
//						-> update
// start->	->check-> ->
//						-> create
//			->check->
//
//			->check->
//
// start stage goroutine parse scv file form URL and line by line transmit to the next stage
//
// check stage it is 5 goroutine which get data from start and check product
// if this product exists (in mongoDb collection) and price changed - transmit to the next stage (update)
// if this product not exist (in mongoDb collection) - transmit to the next stage (create)
//
// create stage get data and insert new product into collection
//
// update stage get data and update product price, updated time and counted of updated price
func (uc *productUC) Fetch(ctx context.Context, url string) error {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	type checkData struct {
		name  string
		price float64
	}

	type updateData struct {
		id    primitive.ObjectID
		price float64
	}

	create := func(ctx context.Context, products <-chan *models.Product, errCh chan error) {
		for p := range products {
			err := uc.productsRepos.Create(ctx, p)
			if err != nil {
				errCh <- status.Errorf(codes.Internal, err.Error())
				return
			}
		}
	}

	update := func(ctx context.Context, inData <-chan *updateData, errCh chan error) {
		for d := range inData {
			err := uc.productsRepos.UpdatePrice(ctx, d.id, d.price)
			if err != nil {
				errCh <- status.Errorf(codes.Internal, err.Error())
				return
			}
		}
	}

	check := func(ctx context.Context, inData <-chan *checkData, errCh chan error) (<-chan *models.Product, <-chan *updateData) {

		createChan := make(chan *models.Product)
		updateChan := make(chan *updateData)

		var wg sync.WaitGroup

		wg.Add(5)

		for i := 0; i < 5; i++ {
			go func() {
				defer wg.Done()
				for d := range inData {
					product, err := uc.productsRepos.Get(ctx, d.name)

					if err != nil {
						if errors.Is(err, models.NotFoundProductError) {
							createChan <- createProduct(d.name, d.price)
							continue
						}
						errCh <- status.Errorf(codes.Internal, err.Error())
						return
					}
					if product.Price != d.price {
						updateChan <- &updateData{id: product.Id, price: d.price}
					}
				}
			}()
		}

		go func() {
			wg.Wait()
			close(createChan)
			close(updateChan)
		}()

		return createChan, updateChan
	}

	start := func(ctx context.Context, resBody io.ReadCloser, errCh chan error) <-chan *checkData {
		checkChan := make(chan *checkData)

		reader := csv.NewReader(resBody)

		go func() {
			defer close(checkChan)
			defer resBody.Close()

			for {
				line, err := reader.Read()
				if err == io.EOF {
					break
				}
				name := line[0]
				price, err := strconv.ParseFloat(line[1], 64)
				if err != nil {
					errCh <- status.Errorf(codes.InvalidArgument, err.Error())
					return
				}

				select {
				case checkChan <- &checkData{name: name, price: price}:
				case <-ctx.Done():
					return
				}
			}
		}()
		return checkChan
	}

	// get external csv file by url
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)

	// start goroutine which parse csv file line by line and send in to data channel
	dataCh := start(ctx, resp.Body, errCh)

	// check run 5 goroutines which receive (name and price) from dataCh chan and check
	// if product exist in the database and after check if the price has changed
	// if not send data to createCh
	// if yes send data to updateCh
	createCh, updateCh := check(ctx, dataCh, errCh)

	// run 2 goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	// run goroutine which receive data from createCh and create new product
	go func() {
		defer wg.Done()
		create(ctx, createCh, errCh)
	}()

	// run goroutine which receive data from updateCh and update existing product
	go func() {
		defer wg.Done()
		update(ctx, updateCh, errCh)
	}()

	wg.Wait()
	close(errCh)

	fmt.Println("num of goroutine: ", runtime.NumGoroutine())
	// check if somewhere have error
	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

//func (uc *productUC) Fetch(ctx context.Context, url string) error {
//
//	// get external csv file by url
//	resp, err := http.Get(url)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//
//	start := time.Now()
//	// create csv reader
//	reader := csv.NewReader(resp.Body)
//
//	for {
//		line, err := reader.Read()
//		if err == io.EOF {
//			break
//		}
//
//		if line[0] == "name" && line[1] == "price" {
//			continue
//		}
//
//		name := line[0]
//		price, err := strconv.ParseFloat(line[1], 64)
//		if err != nil {
//			return err
//		}
//
//		product, getErr := uc.productsRepos.Get(ctx, name)
//		if getErr != nil {
//			if errors.Is(getErr, models.NotFoundProductError) {
//				if err := uc.productsRepos.Create(ctx, createProduct(name, price)); err != nil {
//					return err
//				}
//				continue
//			}
//			return getErr
//		}
//
//		fmt.Println(product.Name, product.Price, price)
//		if product.Price != price {
//			// if prices not the same, increase price updates counter by 1
//			if err := uc.productsRepos.UpdatePrice(ctx, product.Id, price); err != nil {
//				return err
//			}
//		}
//	}
//
//	fmt.Println(time.Since(start))
//	return nil
//}

// newProductUC - return pointer of productUC
func newProductUC(repos repository.ProductsReposInterface) *productUC {
	return &productUC{productsRepos: repos}
}

// createProduct - take name and price
// define all fields for models.Product
// return - pointer for this product
func createProduct(name string, price float64) *models.Product {
	return &models.Product{
		Id:           primitive.NewObjectID(),
		Name:         name,
		Price:        price,
		Updated:      time.Now(),
		PriceUpdates: 0,
	}
}
