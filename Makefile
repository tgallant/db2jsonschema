PACKAGE_NAME=db2jsonschema
CMD_PATH=cmd/db2jsonschema/main.go
BUILD_PATH=/go/src/github.com/tgallant/db2jsonschema
BUILD_FILE=Dockerfile
CI_BUILD_FILE=Dockerfile.ci
CI_IMAGE_TAG=db2j_ci

.PHONY: test

deps:
	go get -v ./...

build:
	go build -v -o ${PACKAGE_NAME} ${CMD_PATH}

build_image:
	./scripts/build_image.sh ${BUILD_FILE} ${PACKAGE_NAME}

build_ci_image:
	./scripts/build_image.sh ${CI_BUILD_FILE} ${CI_IMAGE_TAG}

run_ci_image:
	docker run -v ${PWD}:${BUILD_PATH} ${CI_IMAGE_TAG}

ci: build_ci_image run_ci_image

test:
	go test -v ./...

lint: deps
	./scripts/lint.sh

shellcheck:
	./scripts/shellcheck.sh

test_all: test lint shellcheck

run_actions:
	./scripts/actions.sh ${CI_IMAGE_TAG}

actions: build_ci_image run_actions

clean:
	go clean
	rm ${BINARY_NAME}
