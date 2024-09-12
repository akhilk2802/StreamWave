package grpc

import (
	"backend/internal/stream-ingest/services"
	"backend/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

// StartGRPCServer starts the gRPC server and listens on the specified port
func StartGRPCServer() {
	// Listen on a TCP port
	lis, err := net.Listen("tcp", ":50051") // The gRPC server will run on port 50051
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	// Create a new gRPC server instance
	grpcServer := grpc.NewServer()

	// Register the StreamService with the gRPC server
	proto.RegisterStreamServiceServer(grpcServer, &services.StreamService{})

	log.Println("gRPC server is running on port 50051")

	// Start serving requests
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
