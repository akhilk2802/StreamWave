package main

import (
	localgrpc "backend/internal/stream-processing/grpc"
	"backend/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen on port 50052: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterStreamProcessingServiceServer(grpcServer, &localgrpc.StreamProcessingServer{})

	log.Println("Stream Processing Service running on port 50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
