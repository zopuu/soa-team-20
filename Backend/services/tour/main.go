package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"tour.xws.com/handler"
	"tour.xws.com/repository"
	"tour.xws.com/service"
)

type MongoCollections struct {
	Tours     *mongo.Collection
	KeyPoints *mongo.Collection
}

func initMongoDB() MongoCollections {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27018"
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("tourdb")

	collections := MongoCollections{
		Tours:     db.Collection("tours"),
		KeyPoints: db.Collection("keyPoints"),
	}

	return collections
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Adjust origin as needed; use "*" only if you don't use credentials
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Vary", "Origin") // good practice for caches

		// Allow browsers to send these in preflight
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// If you use cookies/Authorization header and need them on the browser:
		// w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight quickly
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func startServer(tourHandler *handler.TourHandler, keyPointHandler *handler.KeyPointHandler) {
	router := mux.NewRouter().StrictSlash(true)

	//TOUR ENDPOINTS
	router.HandleFunc("/tours", tourHandler.GetAll).Methods("GET")
	router.HandleFunc("/tours/users/{userId}", tourHandler.GetAllByAuthor).Methods("GET")
	router.HandleFunc("/tours/{id}", tourHandler.GetById).Methods("GET")
	router.HandleFunc("/tours", tourHandler.Create).Methods("POST")
	router.HandleFunc("/tours/{id}", tourHandler.Delete).Methods("DELETE")
	router.HandleFunc("/tours/{id}", tourHandler.Update).Methods("PUT")

	//KEYPOINT ENDPOINTS
	router.HandleFunc("/keyPoints", keyPointHandler.GetAll).Methods("GET")
	router.HandleFunc("/keyPoints/tours/{tourId}", keyPointHandler.GetAllByTour).Methods("GET")
	router.HandleFunc("/keyPoints/tours/{tourId}/sortedByCreatedAt", keyPointHandler.GetAllByTourSortedByCreatedAt).Methods("GET")
	router.HandleFunc("/keyPoints", keyPointHandler.Create).Methods("POST")
	router.HandleFunc("/keyPoints/{id}", keyPointHandler.Delete).Methods("DELETE")
	router.HandleFunc("/keyPoints/{id}", keyPointHandler.Update).Methods("PUT")
	router.HandleFunc("/keyPoints/{id}/image", keyPointHandler.GetImage).Methods("GET")

	router.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	handlerWithCors := cors(router)

	println("Server starting")
	log.Fatal(http.ListenAndServe(":8080", handlerWithCors))
}

func main() {
	collections := initMongoDB()

	//TOUR
	tourRepository := &repository.TourRepository{Collection: collections.Tours}
	tourService := &service.TourService{TourRepository: tourRepository}
	tourHandler := &handler.TourHandler{TourService: tourService}
	//KEYPOINT
	keyPointRepository := &repository.KeyPointRepository{Collection: collections.KeyPoints}
	keyPointService := &service.KeyPointService{KeyPointRepository: keyPointRepository}
	keyPointHandler := &handler.KeyPointHandler{KeyPointService: keyPointService}

	startServer(tourHandler, keyPointHandler)
}
