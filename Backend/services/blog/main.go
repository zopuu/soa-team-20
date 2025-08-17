package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"blog.xws.com/handler"
	"blog.xws.com/repository"
	"blog.xws.com/service"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCollections struct {
	Blogs    *mongo.Collection
	Comments *mongo.Collection
}

func initMongoDB() MongoCollections {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("blogdb")

	collections := MongoCollections{
		Blogs:    db.Collection("blogs"),
		Comments: db.Collection("comments"),
	}

	return collections
}

func startServer(handler *handler.BlogHandler) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/blogs", handler.GetAll).Methods("GET")
	router.HandleFunc("/blogs/{id}", handler.GetById).Methods("GET")
	router.HandleFunc("/blogs", handler.Create).Methods("POST")
	router.HandleFunc("/blogs/{id}", handler.Delete).Methods("DELETE")
	router.HandleFunc("/blogs/{id}", handler.Update).Methods("PUT")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	println("Server starting")
	log.Fatal(http.ListenAndServe(":8080", router))
}
func main() {
	collections := initMongoDB()

	blogRepository := &repository.BlogRepository{Collection: collections.Blogs}
	blogService := &service.BlogService{BlogRepository: blogRepository}
	blogHandler := &handler.BlogHandler{BlogService: blogService}

	startServer(blogHandler)
}
