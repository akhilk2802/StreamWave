package grpc

import (
	"/backend/proto"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

// StreamServer implements the StreamService defined in the proto file
type StreamServer struct {
	proto.UnimplementedStreamServiceServer
}

// StartStream handles starting a new stream
func (s *StreamServer) StartStream(ctx context.Context, req *proto.StreamRequest) (*proto.StreamResponse, error) {
	log.Printf("Starting stream: %s for user: %s", req.StreamKey, req.UserId)

	// Implement your business logic for starting the stream
	// e.g., Forward to NGINX RTMP or notify Stream Processing Service

	return &proto.StreamResponse{
		Status:  "success",
		Message: "Stream started successfully",
	}, nil
}

// StopStream handles stopping an active stream
func (s *StreamServer) StopStream(ctx context.Context, req *proto.StreamRequest) (*proto.StreamResponse, error) {
	log.Printf("Stopping stream: %s", req.StreamKey)

	// Implement your business logic for stopping the stream
	// e.g., Clean up resources, stop forwarding to processing service

	return &proto.StreamResponse{
		Status:  "success",
		Message: "Stream stopped successfully",
	}, nil
}

// ForwardMetadata forwards metadata related to the stream
func (s *StreamServer) ForwardMetadata(ctx context.Context, req *proto.MetadataRequest) (*proto.MetadataResponse, error) {
	log.Printf("Forwarding metadata for stream: %s, Metadata: %s", req.StreamKey, req.Metadata)

	// Implement your business logic for forwarding metadata
	// e.g., Send metadata to the Stream Processing Service

	return &proto.MetadataResponse{
		Status:  "success",
		Message: "Metadata forwarded successfully",
	}, nil
}

// StartGRPCServer starts the gRPC server on port 50051
func StartGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterStreamServiceServer(grpcServer, &StreamServer{})

	log.Println("gRPC server running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
