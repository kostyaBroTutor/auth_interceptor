FROM golang:1.20-alpine3.16

WORKDIR /workspace
ENV PATH="${PATH}:/usr/local/go/bin:/root/go/bin:$(go env GOPATH)/bin"

RUN apk add git && \
    git clone https://github.com/protocolbuffers/protobuf.git && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 && \
    go install github.com/vektra/mockery/v2@latest && \
    apk add protoc
