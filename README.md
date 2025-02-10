# StreamWave: Scalable Live Streaming Platform

## Overview
StreamWave is a distributed live streaming platform designed to handle real-time video ingestion, transcoding, and delivery. It utilizes **RTMP** for video ingestion, **FFmpeg** for transcoding, and **DASH** for adaptive streaming. The system is built with **Golang**, leveraging **gRPC** for inter-service communication and designed to scale with **Docker** and **Kubernetes** for cloud-native deployment.

---

## Features
- **Real-time Video Streaming**: Supports live video streaming via RTMP.
- **Automated Video Transcoding**: Converts incoming streams into multiple resolutions (1080p, 720p, 480p, 360p) for adaptive streaming.
- **Segmented Video Processing**: Implements **MapReduce** for efficient segmentation and distributed processing.
- **Cloud Storage Integration**: Supports both local file system storage and AWS S3 for scalable storage.
- **DASH Playback Support**: Outputs video segments in **DASH format** for smooth playback across devices.
- **Microservice Architecture**: Two independent services (**Stream-Ingest** & **Stream-Processing**) communicating via **gRPC**.
<!-- - **Kubernetes & Docker Support**: Fully containerized for seamless cloud deployment and orchestration. -->

---

## Architecture
StreamWave follows a **microservices-based architecture**, consisting of the following services:

### 1. **Stream-Ingest Service** (RTMP Ingestion & Metadata Handling)
- **Listens for live streams** from OBS Studio via RTMP.
- **Triggers metadata processing hooks** (on_publish, on_done) when a stream starts or stops.
- **Communicates with Stream-Processing Service** via gRPC to initiate transcoding.

### 2. **Stream-Processing Service** (Transcoding & Segmentation)
- **Receives a request from Stream-Ingest** to process a new stream.
- **Uses FFmpeg** to transcode video into multiple resolutions.
- **Segments video into chunks** for DASH streaming.
- **Stores output files in the local system or AWS S3**.

---

## Folder Structure
```plaintext
StreamWave/
├── backend/
│   ├── cmd/
│   │   ├── stream_ingest/
│   │   │   └── main.go      # Entry point for Stream-Ingest Service
│   │   ├── stream_processing/
│   │   │   └── main.go      # Entry point for Stream-Processing Service
│   ├── internal/
│   │   ├── stream_ingest/
│   │   │   ├── grpc/
│   │   │   │   └── client.go
│   │   │   ├── handlers/
│   │   │   │   └── stream_ingest_handlers.go
│   │   ├── stream_processing/
│   │   │   ├── mapReduce/
│   │   │   │   └── mapReduce_framework.go
│   │   │   ├── processor/
│   │   │   │   └── segment_processor.go
│   │   │   ├── grpc/
│   │   │   │   └── server.go
│   ├── proto/
│   │   ├── stream_ingest.proto
│   │   ├── stream_processing.proto
│   │   ├── common.proto
├── frontend/   # Yet to complete
├── docs/       # Documentation
├── .gitignore
└── README.md
```

---

## Technologies Used
| Component          | Technology |
|-------------------|------------|
| Backend Framework | Golang (Gin, gRPC) |
| Video Ingestion   | NGINX with RTMP Module |
| Video Transcoding | FFmpeg |
| Video Streaming   | DASH (Dynamic Adaptive Streaming over HTTP) |
| Storage           | Local FS / AWS S3 |
| Communication     | gRPC |
---

## Installation & Setup
### Prerequisites
- **Go** (>=1.18)
- **FFmpeg**
- **NGINX with RTMP Module**
- **Docker & Kubernetes (Optional for cloud deployment)**

### Step 1: Clone the Repository
```sh
git clone https://github.com/akhilk2802/StreamWave.git
cd StreamWave/backend
```

### Step 2: Configure Environment Variables
Create a `.env` file in the backend directory:
```ini
USE_S3=false
OUTPUT_DIR=./output
FFMPEG_PATH=/usr/bin/ffmpeg
RTMP_URL=rtmp://localhost:1936/live
```

### Step 3: Start Services
```sh
# Start Stream-Ingest Service
cd cmd/stream_ingest
go run main.go

# Start Stream-Processing Service
cd ../stream_processing
go run main.go
```

### Step 4: Start NGINX RTMP Server
Ensure your `nginx.conf` is properly set up and run:
```sh
sudo nginx -t
sudo nginx -s reload
```

### Step 5: Test with OBS Studio
- **Set OBS output to RTMP**: `rtmp://localhost:1936/live`
- **Start streaming**
- **Monitor logs in the backend services**

---

## API Endpoints
### Stream-Ingest Service (Port 8081)
| Method | Endpoint           | Description |
|--------|------------------|-------------|
| `POST` | `/start-stream`  | Initiates streaming |
| `POST` | `/stop-stream`   | Stops a stream |
| `GET`  | `/status`        | Checks service health |

### Stream-Processing Service (Port 50052, gRPC)
| Method          | Description |
|----------------|-------------|
| `StartStream()` | Initiates transcoding |
| `StopStream()`  | Stops transcoding |
| `ForwardMetadata()` | Forwards metadata |

---

## Deployment
### **Docker Compose (Local Setup)**
```sh
docker-compose up --build
```

### **Kubernetes (Cloud Deployment)**
```sh
kubectl apply -f k8s/streamwave-deployment.yaml
kubectl get pods -n streamwave
```

---

## Future Improvements
- **Integrate WebRTC for Low-Latency Streaming**
- **Implement Video Analytics for User Engagement**
- **Expand Storage Options with Multi-CDN Support**

---

## Contributors
- **Akhileshkumar S Kumbar** - [GitHub](https://github.com/akhilk2802)

## License
This project is licensed under the **MIT License**.