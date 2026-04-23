#!/bin/bash

set -euo pipefail

start_time=$(date +%s)
start_datetime=$(date "+%Y-%m-%d %H:%M:%S")
echo "Script initiated at: $start_datetime"

##########################################################

# Build a local binary for testing
OUTPUT_BIN="svt-local"

echo "Building local binary: ${OUTPUT_BIN}"
go build -ldflags="-X main.buildMode=local" -o "${OUTPUT_BIN}" ./cmd/ssh/main.go

echo "Build completed successfully"

##########################################################

end_time=$(date +%s)
deployment_time=$((end_time - start_time))
minutes=$((deployment_time / 60))
seconds=$((deployment_time % 60))
echo "Completed in ${minutes}m ${seconds}s"
