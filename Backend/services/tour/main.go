package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	tourgrpc "tour.xws.com/grpc"
	"tour.xws.com/handler"
	tourpb "tour.xws.com/proto"
	"tour.xws.com/repository"
	"tour.xws.com/service"

	obs "github.com/zopuu/soa-team-20/common/obs"
)

type MongoCollections struct {
	Tours            *mongo.Collection
	KeyPoints        *mongo.Collection
	CurrentLocations *mongo.Collection
	Ratings          *mongo.Collection
	TourExecutions   *mongo.Collection
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
		Tours:            db.Collection("tours"),
		KeyPoints:        db.Collection("keyPoints"),
		CurrentLocations: db.Collection("currentLocations"),
		Ratings:          db.Collection("tour_ratings"),
		TourExecutions:   db.Collection("tour_executions"),
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
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID, X-Trace-Id")
		w.Header().Set("Access-Control-Expose-Headers", "X-Request-ID, X-Trace-Id")
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

func startGRPCServer(tourService *service.TourService) {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed listen %s: %v", ":50052", err)
	}

	grpcServer := grpc.NewServer()
	tourGRPCServer := tourgrpc.NewTourGRPCServer(tourService)
	tourpb.RegisterTourServiceServer(grpcServer, tourGRPCServer)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC serve error: %v", err)
	}
}

func startServer(tourHandler *handler.TourHandler, keyPointHandler *handler.KeyPointHandler, locationHandler *handler.CurrentLocationHandler, ratingHandler *handler.TourRatingHandler, tourExecHandler *handler.TourExecutionHandler) {
	router := mux.NewRouter().StrictSlash(true)

	// middleware: logging
	logger := obs.NewLogger("tourservice")
	router.Use(obs.TraceMiddleware)
	router.Use(obs.AccessLogMiddleware(logger))

	//TOUR ENDPOINTS
	router.HandleFunc("/tours", tourHandler.GetAll).Methods("GET")
	router.HandleFunc("/tours/users/{userId}", tourHandler.GetAllByAuthor).Methods("GET")
	router.HandleFunc("/tours/{id}", tourHandler.GetById).Methods("GET")
	router.HandleFunc("/tours", tourHandler.Create).Methods("POST")
	router.HandleFunc("/tours/{id}", tourHandler.Delete).Methods("DELETE")
	router.HandleFunc("/tours/{id}", tourHandler.Update).Methods("PUT")

	router.HandleFunc("/tours/{id}/reviews", ratingHandler.Create).Methods("POST")
	router.HandleFunc("/tours/{id}/reviews", ratingHandler.GetByTour).Methods("GET")
	//KEYPOINT ENDPOINTS
	router.HandleFunc("/keyPoints", keyPointHandler.GetAll).Methods("GET")
	router.HandleFunc("/keyPoints/tours/{tourId}", keyPointHandler.GetAllByTour).Methods("GET")
	router.HandleFunc("/keyPoints/tours/{tourId}/sortedByCreatedAt", keyPointHandler.GetAllByTourSortedByCreatedAt).Methods("GET")
	router.HandleFunc("/keyPoints", keyPointHandler.Create).Methods("POST")
	router.HandleFunc("/keyPoints/{id}", keyPointHandler.Delete).Methods("DELETE")
	router.HandleFunc("/keyPoints/{id}", keyPointHandler.Update).Methods("PUT")
	router.HandleFunc("/keyPoints/{id}/image", keyPointHandler.GetImage).Methods("GET")

	//CURRENT LOCATION ENDPOINTS
	router.HandleFunc("/simulator/location/{userId}", locationHandler.Get).Methods("GET")
	router.HandleFunc("/simulator/location", locationHandler.Set).Methods("PUT")

	router.HandleFunc("/tours/tour-executions/start",           tourExecHandler.Start).Methods("POST")
	router.HandleFunc("/tours/tour-executions/check",           tourExecHandler.CheckProximity).Methods("POST")
	router.HandleFunc("/tours/tour-executions/abandon",         tourExecHandler.Abandon).Methods("POST")
	router.HandleFunc("/tours/tour-executions/active",          tourExecHandler.GetActive).Methods("POST") // <-- was GET
	router.HandleFunc("/tours/tour-executions/active-for-tour", tourExecHandler.GetActiveForTour).Methods("POST")



	router.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	handlerWithCors := cors(router)

	logger.Info("starting_server", obs.F("addr", ":8080"))
	if err := http.ListenAndServe(":8080", handlerWithCors); err != nil {
		logger.Fatal("server_exit", obs.Err(err))
	}
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
	//CURRENT LOCATION
	locationRepo := &repository.CurrentLocationRepository{Collection: collections.CurrentLocations}
	locationService := &service.CurrentLocationService{Repo: locationRepo}
	locationHandler := &handler.CurrentLocationHandler{Svc: locationService}
	//RATINGS
	ratingRepo := &repository.TourRatingRepository{Collection: collections.Ratings}
	ratingService := &service.TourRatingService{Repo: ratingRepo}
	ratingHandler := &handler.TourRatingHandler{RatingService: ratingService}
	// TOUR EXECUTION
	tourExecRepo := &repository.TourExecutionRepository{Collection: collections.TourExecutions}
	tourExecService := &service.TourExecutionService{
		TourRepo:       tourRepository,     // already created above
		KeyPointRepo:   keyPointRepository, // already created above
		CurrentLocRepo: locationRepo,       // already created above
		TourExecRepo:   tourExecRepo,       // new
	}
	tourExecHandler := &handler.TourExecutionHandler{Svc: tourExecService}
	go func() {
        log.Println("Starting gRPC server on port 50052...")
        startGRPCServer(tourService)
    }()

	startServer(tourHandler, keyPointHandler, locationHandler, ratingHandler, tourExecHandler)
	
	
}
