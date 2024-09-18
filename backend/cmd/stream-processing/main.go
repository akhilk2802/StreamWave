package main

import (
	localgrpc "backend/internal/stream-processing/grpc"
	"backend/proto"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("STREAM_PROCESSING_PORT")

	lis, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterStreamProcessingServiceServer(grpcServer, &localgrpc.StreamProcessingServer{})

	log.Printf("Stream Processing Service running on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
