package grpc

import (
    "context"
    "log"

    "backend/proto"
    "google.golang.org/grpc"
)

func StartStreamProcessing(streamKey string) error {
    conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
    if err != nil {
        return err
    }
    defer conn.Close()

    client := proto.NewStreamProcessingServiceClient(conn)
    req := &proto.ProcessingRequest{
        StreamKey: streamKey,
    }

    _, err = client.StartProcessing(context.Background(), req)
    if err != nil {
        log.Printf("Error calling Stream Processing Service: %v", err)
        return err
    }

    log.Println("Successfully started stream processing")
    return nil
}