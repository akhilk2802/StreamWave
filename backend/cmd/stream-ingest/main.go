package streamingest

import (
	"backend/internal/stream-ingest/grpc"
)

func main() {
	// Start the gRPC server
	grpc.StartGRPCServer()

	// Optionally run other services (e.g., HTTP server) here
}
