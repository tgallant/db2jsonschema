BINARY_NAME=db2jsonschema
MAIN_PATH=cmd/db2jsonschema/main.go

.PHONY: build test lint local_ci clean

build:
	go build -v -o ${BINARY_NAME} ${MAIN_PATH}

test:
	go test -v ./...

lint:
	./scripts/lint.sh

test_all: lint test

local_ci:
	./scripts/local_ci.sh

clean:
	go clean
	rm ${BINARY_NAME}
