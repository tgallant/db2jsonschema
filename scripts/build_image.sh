#!/bin/sh

DOCKERFILE=$1
TAG=${2-latest}

docker_exists() {
    type docker > /dev/null 2> /dev/null
}

if ! docker_exists; then
    echo "Docker is not installed. Please install it on your system to continue."
    exit 1
fi

docker build -f "$DOCKERFILE" -t "$TAG" .
