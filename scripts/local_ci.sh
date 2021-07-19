#!/bin/sh

BUILD_TAG=db2j_ci

act_exists() {
    type act > /dev/null 2> /dev/null
}

if ! act_exists; then
    echo "act is not installed. Please install it on your system to continue."
    exit 1
fi

act -P ubuntu-20.04=$BUILD_TAG
