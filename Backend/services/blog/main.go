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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCollections struct {
	Blogs    *mongo.Collection
	Comments *mongo.Collection
	Likes    *mongo.Collection
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
		Likes:    db.Collection("likes"),
	}

	likeIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "userId", Value: 1},
			{Key: "blogId", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err = collections.Likes.Indexes().CreateOne(ctx, likeIndex)
	if err != nil {
		log.Fatalf("Failed to create index on Likes: %v", err)
	}

	return collections
}

func startServer(handler *handler.BlogHandler, commentHandler *handler.CommentHandler, likeHandler *handler.LikeHandler) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/blogs", handler.GetAll).Methods("GET")
	router.HandleFunc("/blogs/{id}", handler.GetById).Methods("GET")
	router.HandleFunc("/blogs", handler.Create).Methods("POST")
	router.HandleFunc("/blogs/{id}", handler.Delete).Methods("DELETE")
	router.HandleFunc("/blogs/{id}", handler.Update).Methods("PUT")
	router.HandleFunc("/blogs/users/{userId}", handler.GetAllByUser).Methods("GET")
	router.HandleFunc("/blogs/{id}/comments", commentHandler.GetByBlogId).Methods("GET")
	router.HandleFunc("/blogs/comments", commentHandler.Create).Methods("POST")
	router.HandleFunc("/blogs/comments/{id}", commentHandler.Update).Methods("PUT")
	router.HandleFunc("/blogs/comments/{id}", commentHandler.Delete).Methods("DELETE")
	router.HandleFunc("/blogs/comments/{id}", commentHandler.GetById).Methods("GET")
	router.HandleFunc("/blogs/likes/{blogId}", likeHandler.GetByBlogId).Methods("GET")
	router.HandleFunc("/blogs/likes", likeHandler.Create).Methods("POST")
	router.HandleFunc("/blogs/likes/{userId}/{blogId}", likeHandler.Delete).Methods("DELETE")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	println("Server starting")
	log.Fatal(http.ListenAndServe(":8080", router))
}
func main() {
	collections := initMongoDB()

	blogRepository := &repository.BlogRepository{Collection: collections.Blogs}
	blogService := &service.BlogService{BlogRepository: blogRepository}
	blogHandler := &handler.BlogHandler{BlogService: blogService}

	commentRepository := &repository.CommentRepository{Collection: collections.Comments}
	commentService := &service.CommentService{CommentRepository: commentRepository}
	commentHandler := &handler.CommentHandler{CommentService: commentService}

	likeRepository := &repository.LikeRepository{Collection: collections.Likes}
	likeService := &service.LikeService{LikeRepository: likeRepository}
	likeHandler := &handler.LikeHandler{LikeService: likeService}

	startServer(blogHandler, commentHandler, likeHandler)
}
