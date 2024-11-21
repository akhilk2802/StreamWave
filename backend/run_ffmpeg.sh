#!/bin/bash

# Variables
STREAM_KEY=$1
RESOLUTION=$2
OUTPUT_FORMAT=$3
OUTPUT_PATH=$4

# Default values (optional)
: ${RESOLUTION:=1920x1080}
: ${OUTPUT_FORMAT:=dash}
: ${OUTPUT_PATH:=./output}

# FFmpeg command
ffmpeg -loglevel debug -report -nostdin -i rtmp://localhost:1936/live/$STREAM_KEY \
  -c:v libx264 -s $RESOLUTION -f $OUTPUT_FORMAT $OUTPUT_PATH/stream.mpd