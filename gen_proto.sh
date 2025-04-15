#!/bin/bash

# Define directories
PROTO_DIR="./api"
OUT_DIR="./api"

# Ensure output directory exists
mkdir -p "$OUT_DIR"

# Generate Go code from proto files
protoc -I="$PROTO_DIR" \
  --go_out=paths=source_relative:"$OUT_DIR" \
  --go-grpc_out=paths=source_relative:"$OUT_DIR" \
  "$PROTO_DIR/permission.proto"

# Print success message
echo "Proto files have been generated successfully."