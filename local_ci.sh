#!/bin/sh

BUILD_TAG=db2j_localci
DOCKERFILE=Dockerfile.localci

docker_exists() {
    type docker > /dev/null 2> /dev/null
}

act_exists() {
    type act > /dev/null 2> /dev/null
}

if ! docker_exists; then
    echo "Docker is not installed. Please install it on your system to continue."
    exit 1
fi

if ! act_exists; then
    echo "act is not installed. Please install it on your system to continue."
    exit 1
fi

docker build -t $BUILD_TAG -f $DOCKERFILE .
act -P ubuntu-20.04=$BUILD_TAG
