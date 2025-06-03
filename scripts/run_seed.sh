#!/bin/sh

# This script runs the Go file data/run_seed.go

# Navigate to the root directory of the project
cd "$(dirname "$0")/.."

# Run the Go file
go run data/run_seed.go