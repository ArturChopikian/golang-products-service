package main

import (
	"fmt"
	csv_server "github.com/ArturChopikian/csv_http_server"
	"github.com/ArturChopikian/grpc-server/configs"
	"github.com/ArturChopikian/grpc-server/database"
	"github.com/ArturChopikian/grpc-server/internal/server"
	"github.com/joho/godotenv"
	"io"
	"log"
	"os"
)

func init() {
	log.SetPrefix("Server: ")

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	f, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	// write logs to file and console
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	cfg, err := configs.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	coll, err := database.NewMongoDBCollection()
	if err != nil {
		log.Fatal(err)
	}

	folder := cfg.SeverCSV.Folder
	address := fmt.Sprintf("%s:%s", cfg.SeverCSV.Host, cfg.SeverCSV.Port)

	csvLog := log.New(wrt, "CSV SERVER: ", log.Ldate|log.Ltime)
	csvServer, err := csv_server.NewCSVServer(folder, address)
	if err != nil {
		csvLog.Fatal(err)
	}

	// map csv handlers and run csv server
	go func() {
		if err := csvServer.MapHandlers(); err != nil {
			csvLog.Fatal(err)
		}
		if err := csvServer.Run(); err != nil {
			csvLog.Fatal(err)
		}
	}()

	// create new products server
	s, err := server.NewProductsServer(cfg, coll)
	if err != nil {
		log.Fatal(err)
	}

	s.MapHandler()

	// run server
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal, 1)

	<-quit

	s.Stop()
}
