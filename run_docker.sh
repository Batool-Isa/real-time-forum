#!/bin/bash

# Define the image name
IMAGE_NAME="my-forum-app"

# Build the Docker image
echo "Building the Docker image..."
docker build -t $IMAGE_NAME .

# Check if the build was successful
if [ $? -ne 0 ]; then
    echo "Docker build failed. Exiting..."
    exit 1
fi

# Run the Docker container
echo "Running the Docker container..."
docker run -p 8000:8000 $IMAGE_NAME

# Check if the container is running
if [ $? -ne 0 ]; then
    echo "Failed to run the Docker container. Exiting..."
    exit 1
fi

echo "Docker container is running. Access the app at http://localhost:8000"
