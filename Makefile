BINARY_NAME=db2jsonschema
MAIN_PATH=cmd/db2jsonschema/main.go

.PHONY: build test lint shellcheck local_ci clean

build:
	go build -v -o ${BINARY_NAME} ${MAIN_PATH}

build_image:
	./scripts/build_image.sh Dockerfile db2jsonschema

test:
	go test -v ./...

lint:
	./scripts/lint.sh

shellcheck:
	./scripts/shellcheck.sh

test_all: lint test shellcheck

run_local_ci:
	./scripts/local_ci.sh

local_ci: build_ci_image run_local_ci

clean:
	go clean
	rm ${BINARY_NAME}
