package main

import (
	"log"
	"net"

	grpcserver "backend/internal/stream-processing/grpc"
	"backend/proto"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
		// Proceeding without .env file
	}

	// Listen on port 50052 for gRPC requests
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen on port 50052: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterStreamProcessingServiceServer(grpcServer, &grpcserver.StreamProcessingServer{})

	log.Println("Stream Processing Service running on port 50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
