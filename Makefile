BINARY_NAME=db2jsonschema
MAIN_PATH=cmd/db2jsonschema/main.go

.PHONY: build test clean

build:
	go build -o ${BINARY_NAME} ${MAIN_PATH}

test:
	go test -v ./...

clean:
	go clean
	rm ${BINARY_NAME}
