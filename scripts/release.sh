#!/bin/sh

git_exists() {
    type git > /dev/null 2> /dev/null
}

semver_exists() {
    type semver > /dev/null 2> /dev/null
}

on_main_branch() {
    CURRENT_BRANCH=$(git branch --show-current)
    echo "Current branch: $CURRENT_BRANCH"
    test "$CURRENT_BRANCH" = main
}

tag_and_push() {
    echo "Creating tag: $1"
    git tag "$1"
    git push origin "$1"
}

yn() {
    printf "%s (y/N) " "$*"
    read -r REPLY
    [ "$REPLY" = "y" ] || [ "$REPLY" = "Y" ]
}

if ! git_exists; then
    echo "Git is not installed. Please install it on your system to continue."
    exit 1
fi

if ! semver_exists; then
    echo "semver is not installed. Please install it on your system to continue."
    exit 1
fi

if ! on_main_branch; then
    echo "Not on the main branch. Please change branches to continue."
    exit 1
fi

INCREMENT_LEVEL=$1
CURRENT_VERSION=$(git tag --list --sort=-version:refname "v*" | head -n 1)
NEXT_VERSION=v$(semver "$CURRENT_VERSION" -i "$INCREMENT_LEVEL")

echo "New $INCREMENT_LEVEL version: $NEXT_VERSION"
yn "Create the tag and push to origin?" && tag_and_push "$NEXT_VERSION"
