package processor

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	// log.Printf("FFmpegPath: %s", sp.FFmpegPath)
	// log.Printf("OutputDir: %s", sp.OutputDir)
	// log.Printf("Bitrate: %s", sp.Bitrate)
	// log.Printf("OutputFormat: %s", sp.OutputFormat)
	// log.Printf("UseS3: %v", sp.UseS3)

	fmt.Printf("OS Environment: %v \n", os.Environ())

	conn, err := net.Dial("tcp", "localhost:1936")
	if err != nil {
		log.Fatalf("Failed to connect to RTMP server: %v", err)
	}
	defer conn.Close()
	fmt.Println("Successfully connected to RTMP server")

	// Check if the FFmpeg executable exists
	if _, err := os.Stat(sp.FFmpegPath); os.IsNotExist(err) {
		log.Printf("FFmpeg executable not found at path: %s", sp.FFmpegPath)
		return fmt.Errorf("FFmpeg executable not found at path: %s", sp.FFmpegPath)
	}

	// Check if the output directory exists
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
	log.Printf("streamURL : %s ", streamURL)

	// if !isStreamAvailableWithRetry(streamURL, 5, 2*time.Second) {
	// 	log.Printf("RTMP stream is not available at URL: %s after retries", streamURL)
	// 	return fmt.Errorf("RTMP stream is not available at URL: %s after retries", streamURL)
	// }

	// Construct the FFmpeg command
	// ffmpegCmd := fmt.Sprintf("ffmpeg -nostdin -i rtmp://localhost:1936/live/%s -c:v libx264 -s %s -f %s %s/stream.mpd",
	// 	streamKey, resolution, sp.OutputFormat, outputPath)
	// log.Printf("Executing FFmpeg command: %s", ffmpegCmd)
	// ffmpegCmd := fmt.Sprintf("/opt/homebrew/bin/ffmpeg -nostdin -i rtmp://localhost:1936/live/%s -c:v libx264 -s %s -f %s %s/stream.mpd",
	// 	streamKey, resolution, sp.OutputFormat, outputPath)
	// log.Printf("Executing FFmpeg command: %s", ffmpegCmd)

	ffmpegCommand := "ffmpeg -i rtmp://localhost:1936/live/test -c:v libx264 -s 1920x1080 -f dash ./output/test/1080p/stream.mpd"

	// Prepare the command execution with a timeout context
	// ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second) // Set a 60-second timeout
	// defer cancel()

	// cmd := exec.Command("bash", "-c", ffmpegCmd)
	parts := strings.Fields(ffmpegCommand)
	cmd := exec.Command(parts[0], parts[1:]...)
	// cmd := exec.CommandContext(ctx,
	// 	"ffmpeg",
	// 	"-nostdin",
	// 	"-i", "rtmp://localhost:1936/live/"+streamKey,
	// 	"-c:v", "libx264",
	// 	"-s", resolution,
	// 	"-f", sp.OutputFormat,
	// 	outputPath+"/stream.mpd")
	// cmd.Env = append(os.Environ(), fmt.Sprintf("PATH=%s", os.Getenv("PATH")))

	// Redirect stdout and stderr to the console
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the FFmpeg command
	if err := cmd.Start(); err != nil {
		log.Printf("Error starting transcoding: %v", err)
		return fmt.Errorf("failed to start FFmpeg command: %v", err)
	}

	// Wait for the FFmpeg command to finish
	err = cmd.Wait()
	if err != nil {
		log.Printf("Error during transcoding for resolution %s: %v", resName, err)
		return fmt.Errorf("transcoding failed for resolution %s: %v", resName, err)
	}

	log.Printf("Transcoding completed successfully for resolution: %s", resName)

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

			key := fmt.Sprintf("%s/%s/%s", streamKey, resName, filepath.Base(path))

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
