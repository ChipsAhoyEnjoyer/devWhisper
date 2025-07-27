#!/bin/bash

# Get the current git commit hash
COMMIT_HASH=$(git rev-parse --short HEAD)

# Check if the git command was successful
if [ $? -ne 0 ]; then
    echo "Error: Failed to get git commit hash. Are you in a git repository?"
    exit 1
fi

echo "Building Docker image with tags: $COMMIT_HASH and latest"

# Build the Docker image with both tags
docker build -t devwhisper:$COMMIT_HASH -t devwhisper:latest .

# Check if the build was successful
if [ $? -eq 0 ]; then
    echo "Docker image successfully built with tags:"
    echo "  - devwhisper:$COMMIT_HASH"
    echo "  - devwhisper:latest"
else
    echo "Error: Docker build failed"
    exit 1
fi