BINARY_NAME=db2jsonschema
MAIN_PATH=cmd/db2jsonschema/main.go

.PHONY: build test lint shellcheck local_ci clean

build:
	go build -v -o ${BINARY_NAME} ${MAIN_PATH}

test:
	go test -v ./...

lint:
	./scripts/lint.sh

shellcheck:
	./scripts/shellcheck.sh

test_all: lint test shellcheck

local_ci:
	./scripts/local_ci.sh

clean:
	go clean
	rm ${BINARY_NAME}
