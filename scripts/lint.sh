#!/bin/sh

golangcilint_exists() {
    type golangci-lint > /dev/null 2> /dev/null
}

if ! golangcilint_exists; then
    echo "golangci-lint is not installed. Please install it on your system to continue."
    exit 1
fi

golangci-lint run
