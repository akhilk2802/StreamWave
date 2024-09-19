package processor

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type SegmentProcessor struct {
	FFmpegPath   string
	Bitrate      string
	OutputFormat string
	OutputDir    string
	UseS3        bool
	S3Bucket     string
	S3Region     string
}

// TranscodeVideo transcodes video into different resolutions
// func (sp *SegmentProcessor) TranscodeVideo(streamKey, resName, resolution string) error {
// 	outputPath := fmt.Sprintf("%s/%s/%s", sp.OutputDir, streamKey, resName)
// 	err := os.MkdirAll(outputPath, 0755)
// 	if err != nil {
// 		return fmt.Errorf("failed to create output directory: %v", err)
// 	}

// 	ffmpegCmd := fmt.Sprintf("%s -i rtmp://localhost:1936/live/%s -c:v libx264 -b:v %s -s %s -f %s %s/stream.mpd",
// 		sp.FFmpegPath, streamKey, sp.Bitrate, resolution, sp.OutputFormat, outputPath)

// 	log.Printf("Executing FFmpeg command for resolution %s: %s", resName, ffmpegCmd)

// 	cmd := exec.Command("bash", "-c", ffmpegCmd)
// 	var out, errBuf bytes.Buffer
// 	cmd.Stdout = &out
// 	cmd.Stderr = &errBuf

// 	if err := cmd.Run(); err != nil {
// 		log.Printf("Error during transcoding for resolution %s: %v", resName, err)
// 		return err
// 	}

// 	log.Printf("Transcoding completed for resolution: %s", resName)
// 	return nil
// }

func (sp *SegmentProcessor) TranscodeVideo(streamKey, resName, resolution string) error {
	log.Printf("Starting TranscodeVideo for streamKey: %s, resName: %s, resolution: %s", streamKey, resName, resolution)

	// Log the environment variables and configurations
	log.Printf("FFmpegPath: %s, OutputDir: %s, Bitrate: %s, OutputFormat: %s, UseS3: %v", sp.FFmpegPath, sp.OutputDir, sp.Bitrate, sp.OutputFormat, sp.UseS3)

	// Ensure the output directory exists
	outputPath := fmt.Sprintf("%s/%s/%s", sp.OutputDir, streamKey, resName)
	err := os.MkdirAll(outputPath, 0755)
	if err != nil {
		log.Printf("Failed to create output directory: %v", err)
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Construct the FFmpeg command
	ffmpegCmd := fmt.Sprintf("%s -y -i rtmp://localhost:1936/live/%s -c:v libx264 -b:v %s -s %s -f %s %s/stream.mpd",
		sp.FFmpegPath, streamKey, sp.Bitrate, resolution, sp.OutputFormat, outputPath)

	log.Printf("Executing FFmpeg command for resolution %s: %s", resName, ffmpegCmd)

	// Verify that the FFmpeg executable exists
	if _, err := os.Stat(sp.FFmpegPath); os.IsNotExist(err) {
		log.Printf("FFmpeg executable not found at path: %s", sp.FFmpegPath)
		return fmt.Errorf("FFmpeg executable not found at path: %s", sp.FFmpegPath)
	}

	// Check if the RTMP stream is available
	streamURL := fmt.Sprintf("rtmp://localhost:1936/live/%s", streamKey)
	if !isStreamAvailable(streamURL) {
		log.Printf("RTMP stream is not available at URL: %s", streamURL)
		return fmt.Errorf("RTMP stream is not available at URL: %s", streamURL)
	}

	// Prepare the command execution
	cmd := exec.Command("bash", "-c", ffmpegCmd)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	// Run the FFmpeg command
	if err := cmd.Run(); err != nil {
		log.Printf("Error during transcoding for resolution %s: %v. FFmpeg stderr: %s", resName, err, errBuf.String())
		return fmt.Errorf("transcoding failed for resolution %s: %v", resName, err)
	}

	log.Printf("Transcoding completed successfully for resolution: %s", resName)
	return nil
}

func isStreamAvailable(streamURL string) bool {
	// Use FFprobe to check if the stream is available
	ffprobeCmd := fmt.Sprintf("ffprobe -v error -show_entries stream=codec_type -of default=noprint_wrappers=1 %s", streamURL)
	cmd := exec.Command("bash", "-c", ffprobeCmd)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	if err != nil {
		log.Printf("FFprobe error: %v. FFprobe stderr: %s", err, errBuf.String())
		return false
	}

	// If FFprobe succeeds, the stream is available
	return true
}

// StoreSegmentsInS3 uploads the processed segments to S3
func (sp *SegmentProcessor) StoreSegmentsInS3(streamKey, resName string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(sp.S3Region),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %v", err)
	}

	svc := s3.New(sess)
	localDir := filepath.Join(sp.OutputDir, streamKey, resName)

	err = filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			key := fmt.Sprintf("%s/%s/%s/%s", streamKey, resName, filepath.Base(path))

			_, err = svc.PutObject(&s3.PutObjectInput{
				Bucket: aws.String(sp.S3Bucket),
				Key:    aws.String(key),
				Body:   file,
				ACL:    aws.String("public-read"),
			})
			if err != nil {
				return err
			}
			log.Printf("Uploaded %s to S3 bucket %s", key, sp.S3Bucket)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to upload segments to S3: %v", err)
	}

	return nil
}
