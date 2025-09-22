package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/zopuu/soa-team-20/Backend/services/shopping_service/internal/api"
	shoppingpb "github.com/zopuu/soa-team-20/Backend/services/shopping_service/proto/shoppingpb"
	obs "github.com/zopuu/soa-team-20/common/obs"
)

func mustGetenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	mongoURI := mustGetenv("MONGO_URI", "mongodb://shopping-mongo:27017")
	dbName := mustGetenv("MONGO_DB", "shoppingdb")
	grpcAddr := mustGetenv("GRPC_ADDR", ":50052") // service listens on :50052

	// Connect to mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}
	db := client.Database(dbName)

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	m := obs.NewMetrics("shopping_service")
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			m.GRPCUnary(),
		),
	)
	go func(){ _ = m.ServeMetrics(":2112") }()
	shoppingpb.RegisterShoppingServiceServer(s, api.NewServer(db))

	reflection.Register(s)

	go func() {
		log.Printf("shopping service gRPC listening on %s", grpcAddr)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("grpc serve: %v", err)
		}
	}()

	// graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("shutting down")
	s.GracefulStop()
	_ = client.Disconnect(context.Background())
}
