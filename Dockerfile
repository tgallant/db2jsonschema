FROM golang:1.16.6-alpine3.14 AS builder
WORKDIR $GOPATH/src/github.com/tgallant/db2jsonschema
COPY . .
RUN apk add --update gcc musl-dev make
RUN make build

FROM alpine:3.14
WORKDIR /root/
COPY --from=builder /go/src/github.com/tgallant/db2jsonschema/db2jsonschema .
ENTRYPOINT ["./db2jsonschema"]
CMD ["help"]
