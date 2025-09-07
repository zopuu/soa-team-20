package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/zopuu/soa-team-20/Backend/services/followers_service/internal/db"
	"github.com/zopuu/soa-team-20/Backend/services/followers_service/internal/server"
	"github.com/zopuu/soa-team-20/Backend/services/followers_service/proto/followerspb"
	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"
)

func main() {
	// Connect Neo4j
	driver := db.NewDriver()
	defer db.CloseDriver(context.Background(), driver)

	// gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	followerspb.RegisterFollowersServiceServer(grpcServer, &server.FollowersServer{Driver: driver})
	reflection.Register(grpcServer)

	fmt.Println("Followers service running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
