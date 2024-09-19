package grpc

import (
	"context"
	"log"
	"os"

	"backend/internal/stream-processing/mapreduce"
	"backend/internal/stream-processing/processor"
	"backend/proto"
)

type StreamProcessingServer struct {
	proto.UnimplementedStreamProcessingServiceServer
}

// StartProcessing handles the request to start stream processing
func (s *StreamProcessingServer) StartProcessing(ctx context.Context, req *proto.ProcessingRequest) (*proto.ProcessingResponse, error) {
	log.Printf("Received request to process stream: %s", req.StreamKey)

	useS3 := os.Getenv("USE_S3") == "true"
	s3Bucket := os.Getenv("S3_BUCKET")
	s3Region := os.Getenv("S3_REGION")
	outputDir := os.Getenv("OUTPUT_DIR")
	ffmpegPath := os.Getenv("FFMPEG_PATH")

	if outputDir == "" {
		outputDir = "./output" // Default output directory
	}

	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg" // Default to ffmpeg in PATH
	}

	segmentProcessor := &processor.SegmentProcessor{
		FFmpegPath:   ffmpegPath,
		Bitrate:      "4500k",
		OutputFormat: "dash",
		OutputDir:    outputDir,
		UseS3:        useS3,
		S3Bucket:     s3Bucket,
		S3Region:     s3Region,
	}

	mapReduce := &mapreduce.MapReduceFramework{}

	// Define resolutions for transcoding
	resolutions := map[string]string{
		"1080p": "1920x1080",
		// "720p":  "1280x720",
		// "480p":  "854x480",
		// "360p":  "640x360",
	}

	// Map phase: Transcoding
	err := mapReduce.Map(segmentProcessor, req.StreamKey, resolutions)
	if err != nil {
		log.Printf("Error during map phase: %v", err)
		return &proto.ProcessingResponse{Status: "error", Message: "Failed during map phase"}, err
	}

	// Reduce phase: Store segments
	err = mapReduce.Reduce(segmentProcessor, req.StreamKey, resolutions)
	if err != nil {
		log.Printf("Error during reduce phase: %v", err)
		return &proto.ProcessingResponse{Status: "error", Message: "Failed during reduce phase"}, err
	}

	return &proto.ProcessingResponse{Status: "success", Message: "Stream processing completed successfully"}, nil
}
