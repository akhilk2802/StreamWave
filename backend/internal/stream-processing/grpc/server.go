package grpc

import (
	"backend/proto"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"

	// "fmt"
	"log"
	// "os"
	// "os/exec"
)

// StreamProcessingServer implements the gRPC methods
type StreamProcessingServer struct {
	proto.UnimplementedStreamProcessingServiceServer
}

// // StartTranscoding starts the video transcoding process
// func (s *StreamProcessingServer) StartTranscoding(ctx context.Context, req *proto.TranscodingRequest) (*proto.TranscodingResponse, error) {
// 	log.Printf("Received request to transcode stream: %s", req.StreamKey)

// 	outputDir := os.Getenv("OUTPUT_DIR")

// 	// ffmpegCmd := fmt.Sprintf("ffmpeg -i rtmp://localhost:1936/live/%s -c:v libx264 -s %s -f %s %s/%s/stream.mpd",
// 	// 	req.StreamKey, req.Resolution, req.Format, outputDir, req.StreamKey)

// 	ffmpegCmd := fmt.Sprintf("ffmpeg -i rtmp://localhost:1936/live/%s -c:v libx264 -s %s -f %s %s/stream.mpd",
// 		req.StreamKey, req.Resolution, req.Format, outputDir)

// 	log.Printf("Executing FFmpeg command: %s", ffmpegCmd)

// 	if err := exec.Command("bash", "-c", ffmpegCmd).Run(); err != nil {
// 		log.Printf("Error starting transcoding: %v", err)
// 		return &proto.TranscodingResponse{Status: "error", Message: "Failed to start transcoding"}, err
// 	}

// 	return &proto.TranscodingResponse{Status: "success", Message: "Transcoding started successfully"}, nil
// }

func (s *StreamProcessingServer) StartTranscoding(ctx context.Context, req *proto.TranscodingRequest) (*proto.TranscodingResponse, error) {
	log.Printf("Received request to transcode stream: %s", req.StreamKey)

	// Step 1: Retrieve and check environment variables (e.g., OUTPUT_DIR)
	outputDir := os.Getenv("OUTPUT_DIR")
	if outputDir == "" {
		log.Printf("Error: OUTPUT_DIR environment variable is not set")
		return &proto.TranscodingResponse{Status: "error", Message: "Output directory not set"}, fmt.Errorf("output directory not set")
	}

	// Ensure output directory exists and is writable
	if err := ensureDirectoryExists(outputDir); err != nil {
		log.Printf("Error: Could not create or access output directory: %v", err)
		return &proto.TranscodingResponse{Status: "error", Message: "Failed to access output directory"}, err
	}

	// Step 2: Build the FFmpeg command
	ffmpegCmd := fmt.Sprintf("ffmpeg -i rtmp://localhost:1936/live/%s -c:v libx264 -s %s -f %s %s/stream.mpd",
		req.StreamKey, req.Resolution, req.Format, outputDir)

	log.Printf("Executing FFmpeg command: %s", ffmpegCmd)

	// Step 3: Prepare command execution with output capture for logging
	cmd := exec.Command("bash", "-c", ffmpegCmd)

	// Capture stdout and stderr for logging
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf

	// Step 4: Execute the FFmpeg command asynchronously
	if err := cmd.Start(); err != nil {
		log.Printf("Error starting transcoding: %v. FFmpeg stderr: %s", err, errBuf.String())
		return &proto.TranscodingResponse{Status: "error", Message: "Failed to start transcoding"}, err
	}

	// Step 5: Wait for the process to finish
	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Printf("Transcoding process failed: %v. FFmpeg stderr: %s", err, errBuf.String())
		} else {
			log.Printf("Transcoding completed successfully. FFmpeg stdout: %s", out.String())
		}
	}()

	log.Printf("Transcoding process started successfully")
	return &proto.TranscodingResponse{Status: "success", Message: "Transcoding started successfully"}, nil
}

// ensureDirectoryExists creates the directory if it does not exist, and checks for write permissions
func ensureDirectoryExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Printf("Directory %s does not exist. Creating...", dir)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	// Check write permissions
	testFile := fmt.Sprintf("%s/testfile", dir)
	if f, err := os.Create(testFile); err != nil {
		return fmt.Errorf("no write permission to directory %s: %v", dir, err)
	} else {
		f.Close()
		os.Remove(testFile) // Clean up test file
	}

	return nil
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
