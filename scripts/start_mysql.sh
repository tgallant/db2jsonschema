#!/bin/sh

TAG=${1-8.0}

docker_exists() {
    type docker > /dev/null 2> /dev/null
}

if ! docker_exists; then
    echo "Docker is not installed. Please install it on your system to continue."
    exit 1
fi

docker run \
       -e MYSQL_ROOT_PASSWORD=root \
       -e MYSQL_DATABASE=testing \
       -p 3306:3306 \
       "mysql:$TAG"
