package mapreduce

import (
	"log"
	"sync"

	"backend/internal/stream-processing/processor"
)

type MapReduceFramework struct{}

// Map phase: Distributes transcoding tasks across goroutines
func (mr *MapReduceFramework) Map(segmentProcessor *processor.SegmentProcessor, streamKey string, resolutions map[string]string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(resolutions))

	for resName, resolution := range resolutions {
		wg.Add(1)
		go func(resName, resolution string) {
			defer wg.Done()
			err := segmentProcessor.TranscodeVideo(streamKey, resName, resolution)
			if err != nil {
				log.Printf("Error transcoding %s: %v", resName, err)
				errChan <- err
			}
		}(resName, resolution)
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

// Reduce phase: Stores the segments (either locally or in S3)
func (mr *MapReduceFramework) Reduce(segmentProcessor *processor.SegmentProcessor, streamKey string, resolutions map[string]string) error {
	useS3 := segmentProcessor.UseS3

	if useS3 {
		log.Println("Uploading segments to S3")
		for resName := range resolutions {
			err := segmentProcessor.StoreSegmentsInS3(streamKey, resName)
			if err != nil {
				return err
			}
		}
	} else {
		log.Println("Storing segments locally")
		// Segments are already stored locally after transcoding
	}

	return nil
}
