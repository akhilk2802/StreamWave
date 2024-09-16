package grpc

import (
	"backend/proto"
	"context"
	"log"

	"google.golang.org/grpc"
)

var streamProcessingClient proto.StreamProcessingServiceClient

// Initialize gRPC client for the Stream Processing Service
func InitStreamProcessingClient() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Stream Processing Service: %v", err)
	}
	streamProcessingClient = proto.NewStreamProcessingServiceClient(conn)
}

// Start transcoding by calling the Stream Processing Service
func StartTranscoding(streamKey, resolution, format string) {
	req := &proto.TranscodingRequest{
		StreamKey:  streamKey,
		Resolution: resolution,
		Format:     format,
	}

	resp, err := streamProcessingClient.StartTranscoding(context.Background(), req)
	if err != nil {
		log.Printf("Error starting transcoding: %v", err)
	} else {
		log.Printf("Transcoding started: %v", resp.Message)
	}
}

// Forward metadata to the Stream Processing Service
func ForwardMetadata(streamKey, metadata string) {
	req := &proto.MetadataRequest{
		StreamKey: streamKey,
		Metadata:  metadata,
	}

	resp, err := streamProcessingClient.ReceiveMetadata(context.Background(), req)
	if err != nil {
		log.Printf("Error forwarding metadata: %v", err)
	} else {
		log.Printf("Metadata forwarded: %v", resp.Message)
	}
}
