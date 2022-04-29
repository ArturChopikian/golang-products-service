package main

import (
	"github.com/ArturChopikian/grpc-server/configs"
	"github.com/ArturChopikian/grpc-server/internal/database"
	"github.com/ArturChopikian/grpc-server/internal/server"
	"io"
	"log"
	"os"
)

func init() {
	log.SetPrefix("Server: ")
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
