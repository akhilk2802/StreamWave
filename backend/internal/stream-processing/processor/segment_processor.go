package processor

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

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


func (sp *SegmentProcessor) TranscodeVideo(streamKey, resName, resolution string) error {
	log.Printf("Starting TranscodeVideo for streamKey: %s, resName: %s, resolution: %s", streamKey, resName, resolution)

	// Log the current working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting working directory: %v", err)
	} else {
		log.Printf("Current working directory: %s", wd)
	}

	// Log the environment variables and configurations
	log.Printf("FFmpegPath: %s", sp.FFmpegPath)
	log.Printf("OutputDir: %s", sp.OutputDir)
	log.Printf("Bitrate: %s", sp.Bitrate)
	log.Printf("OutputFormat: %s", sp.OutputFormat)
	log.Printf("UseS3: %v", sp.UseS3)

	// Ensure the FFmpeg executable exists
	if _, err := os.Stat(sp.FFmpegPath); os.IsNotExist(err) {
		log.Printf("FFmpeg executable not found at path: %s", sp.FFmpegPath)
		return fmt.Errorf("FFmpeg executable not found at path: %s", sp.FFmpegPath)
	}

	// Ensure the output directory exists
	outputPath := filepath.Join(sp.OutputDir, streamKey, resName)
	log.Printf("Output path: %s", outputPath)

	err = os.MkdirAll(outputPath, 0755)
	if err != nil {
		log.Printf("Failed to create output directory: %v", err)
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	if !isDirectoryWritable(outputPath) {
		log.Printf("Output directory is not writable: %s", outputPath)
		return fmt.Errorf("output directory is not writable: %s", outputPath)
	}

	// Check if the RTMP stream is available with retries
	streamURL := fmt.Sprintf("rtmp://localhost:1936/live/%s", streamKey)
	if !isStreamAvailableWithRetry(streamURL, 5, 2*time.Second) {
		log.Printf("RTMP stream is not available at URL: %s after retries", streamURL)
		return fmt.Errorf("RTMP stream is not available at URL: %s after retries", streamURL)
	}

	// Construct the FFmpeg command with reconnection options
	// ffmpegCmd := fmt.Sprintf("%s -y -loglevel debug -reconnect 1 -reconnect_streamed 1 -reconnect_delay_max 5 -i %s -c:v libx264 -b:v %s -s %s -f %s %s/stream.mpd",
	// 	sp.FFmpegPath, streamURL, sp.Bitrate, resolution, sp.OutputFormat, outputPath)

	// log.Printf("Executing FFmpeg command for resolution %s: %s", resName, ffmpegCmd)

	ffmpegCmd := fmt.Sprintf("ffmpeg -i rtmp://localhost:1936/live/%s -c:v libx264 -s %s -f %s %s/stream.mpd",
		streamKey, resolution, sp.OutputFormat, outputPath)
	log.Printf("Executing FFmpeg command: %s", ffmpegCmd)

	// Prepare the command execution with a timeout
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", ffmpegCmd)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PATH=%s", os.Getenv("PATH")))

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	if err := cmd.Start(); err != nil {
		log.Printf("Error starting transcoding: %v", err)
		return fmt.Errorf("failed to start FFmpeg command: %v", err)
	}

	// Wait for the FFmpeg command to finish
	err = cmd.Wait()
	if err != nil {
		log.Printf("Error during transcoding for resolution %s: %v", resName, err)
		log.Printf("FFmpeg stderr: %s", errBuf.String())
		log.Printf("FFmpeg stdout: %s", outBuf.String())
		return fmt.Errorf("transcoding failed for resolution %s: %v", resName, err)
	}

	log.Printf("Transcoding completed successfully for resolution: %s", resName)
	log.Printf("FFmpeg stderr: %s", errBuf.String())
	log.Printf("FFmpeg stdout: %s", outBuf.String())

	return nil
}

func isDirectoryWritable(dirPath string) bool {
	testFile := filepath.Join(dirPath, ".writetest")
	file, err := os.Create(testFile)
	if err != nil {
		return false
	}
	file.Close()
	os.Remove(testFile)
	return true
}

// isStreamAvailable checks if the RTMP stream is available using ffprobe with a timeout
func isStreamAvailable(streamURL string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ffprobeCmd := fmt.Sprintf("ffprobe -v error -rw_timeout 5000000 -show_entries stream=codec_type -of default=noprint_wrappers=1 %s", streamURL)
	cmd := exec.CommandContext(ctx, "bash", "-c", ffprobeCmd)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PATH=%s", os.Getenv("PATH")))

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		log.Printf("FFprobe command timed out while checking stream availability for URL: %s", streamURL)
		return false
	}
	if err != nil {
		log.Printf("FFprobe error: %v. FFprobe stderr: %s", err, errBuf.String())
		return false
	}

	// Optionally, parse outBuf to confirm stream contains expected codecs
	return true
}

// isStreamAvailableWithRetry checks if the RTMP stream is available, retrying if necessary
func isStreamAvailableWithRetry(streamURL string, retries int, delay time.Duration) bool {
	for i := 0; i < retries; i++ {
		log.Printf("Checking stream availability (attempt %d/%d)...", i+1, retries)
		if isStreamAvailable(streamURL) {
			log.Printf("Stream is available at URL: %s", streamURL)
			return true
		}
		log.Printf("Stream not available yet, retrying in %v...", delay)
		time.Sleep(delay)
	}
	log.Printf("Stream not available after %d retries.", retries)
	return false
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
