#!/bin/sh

shellcheck_exists() {
    type shellcheck > /dev/null 2> /dev/null
}

if ! shellcheck_exists; then
    echo "shellcheck is not installed. Please install it on your system to continue."
    exit 1
fi

shellcheck ./*.sh scripts/*.sh
