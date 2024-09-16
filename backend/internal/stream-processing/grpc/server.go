package grpc

import (
	"backend/proto"
	"context"
	"fmt"
	"log"
	"os/exec"
)

// StreamProcessingServer implements the gRPC methods
type StreamProcessingServer struct {
	proto.UnimplementedStreamProcessingServiceServer
}

// StartTranscoding starts the video transcoding process
func (s *StreamProcessingServer) StartTranscoding(ctx context.Context, req *proto.TranscodingRequest) (*proto.TranscodingResponse, error) {
	log.Printf("Received request to transcode stream: %s", req.StreamKey)

	ffmpegCmd := fmt.Sprintf("ffmpeg -i rtmp://localhost/live/%s -c:v libx264 -s %s -f %s ./backend/output/%s/stream.mpd",
		req.StreamKey, req.Resolution, req.Format, req.StreamKey)

	if err := exec.Command("bash", "-c", ffmpegCmd).Run(); err != nil {
		log.Printf("Error starting transcoding: %v", err)
		return &proto.TranscodingResponse{Status: "error", Message: "Failed to start transcoding"}, err
	}

	return &proto.TranscodingResponse{Status: "success", Message: "Transcoding started successfully"}, nil
}

// ReceiveMetadata receives forwarded metadata from the Stream Ingest Service
func (s *StreamProcessingServer) ReceiveMetadata(ctx context.Context, req *proto.MetadataRequest) (*proto.MetadataResponse, error) {
	log.Printf("Received metadata for stream %s: %s", req.StreamKey, req.Metadata)

	// Process metadata (e.g., store it, use it for analytics, etc.)
	return &proto.MetadataResponse{
		Status:  "success",
		Message: "Metadata processed successfully",
	}, nil
}
