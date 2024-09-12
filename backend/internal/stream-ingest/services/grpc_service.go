package services

import (
	"backend/proto"
	"context"
	"fmt"
	"log"
	"os/exec"
)

type StreamService struct {
	proto.UnimplementedStreamServiceServer
}

// gRPC method to start the stream
func (s *StreamService) StartStream(ctx context.Context, req *proto.StreamRequest) (*proto.StreamResponse, error) {
	log.Printf("Starting stream: %s for user: %s", req.StreamKey, req.UserId)

	err := transcodeAndSegmentStream(req.StreamKey)
	if err != nil {
		log.Printf("Error transcoding stream: %s", err.Error())
		return &proto.StreamResponse{
			Status:  "failure",
			Message: fmt.Sprintf("Failed to start stream: %v", err),
		}, err
	}

	// Log the success
	log.Printf("Stream %s started successfully", req.StreamKey)
	return &proto.StreamResponse{
		Status:  "success",
		Message: "Stream started successfully",
	}, nil
}

// gRPC method to stop the stream
func (s *StreamService) StopStream(ctx context.Context, req *proto.StreamRequest) (*proto.StreamResponse, error) {
	log.Printf("Stopping stream: %s", req.StreamKey)

	// Implement your business logic for stopping the stream
	// For example, cleanup resources

	return &proto.StreamResponse{
		Status:  "success",
		Message: "Stream stopped successfully",
	}, nil
}

func transcodeAndSegmentStream(streamKey string) error {
	log.Println("StreamKey", streamKey)
	// Define the input and output streams
	rtmpInput := fmt.Sprintf("rtmp://localhost/live/%s", streamKey)     // RTMP input URL from NGINX
	hlsOutput := fmt.Sprintf("./backend/hls/%s/output.m3u8", streamKey) // HLS output

	log.Printf("rtmpInput: %s", rtmpInput)
	log.Printf("hlsoutput: %s", hlsOutput)
	// Run FFmpeg to transcode and segment the stream
	cmd := exec.Command("ffmpeg",
		"-i", rtmpInput, // Input stream
		"-c:v", "libx264", // Video codec
		"-crf", "23", // Quality factor
		"-preset", "fast", // Encoding preset
		"-f", "hls", // Output format: HLS
		"-hls_time", "10", // Segment duration: 10 seconds
		"-hls_playlist_type", "event", // HLS playlist type
		"-hls_segment_filename", fmt.Sprintf("./backend/hls/%s/segment_%%03d.ts", streamKey), // HLS segment files
		hlsOutput, // Output playlist (.m3u8)
	)

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start transcoding: %w", err)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("transcoding process failed: %w", err)
	}

	log.Printf("Stream %s successfully transcoded and segmented", streamKey)
	return nil
}
