package services

import (
	"fmt"
	"log"
	"os/exec"
)

// TranscodeStream runs the FFmpeg command to transcode the stream
func TranscodeStream(streamKey, resolution, format string) error {
	ffmpegCmd := fmt.Sprintf("ffmpeg -i rtmp://localhost/live/%s -c:v libx264 -s %s -f %s ./backend/output/%s/stream.mpd",
		streamKey, resolution, format, streamKey)

	cmd := exec.Command("bash", "-c", ffmpegCmd)
	if err := cmd.Run(); err != nil {
		log.Printf("Error starting transcoding: %v", err)
		return err
	}

	log.Printf("Transcoding started for stream: %s", streamKey)
	return nil
}
