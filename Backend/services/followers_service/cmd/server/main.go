package main

import (
	"context"
	"log"
	"net"

	"github.com/zopuu/soa-team-20/Backend/services/followers_service/internal/db"
	"github.com/zopuu/soa-team-20/Backend/services/followers_service/internal/server"
	"github.com/zopuu/soa-team-20/Backend/services/followers_service/proto/followerspb"
	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"

	obs "github.com/zopuu/soa-team-20/common/obs"
)

func main() {
	logger := obs.NewLogger("followers")
	// Connect Neo4j
	driver := db.NewDriver()
	defer db.CloseDriver(context.Background(), driver)

	// gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Fatal("listen_error", obs.Err(err))
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			obs.GRPCTraceUnary(),
			obs.GRPCAccessLogUnary(logger),
		),
	)
	followerspb.RegisterFollowersServiceServer(grpcServer, &server.FollowersServer{Driver: driver})
	reflection.Register(grpcServer)

	logger.Info("starting_grpc", obs.F("addr", ":50051"))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("serve_error", obs.Err(err))
	}
}
