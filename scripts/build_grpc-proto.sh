#!/bin/sh
echo "running git submodule update --init --recursive"
git submodule update --init --recursive
echo "running protoc for grpc-proto/grpc/reflection/v1/reflection.proto"
mkdir -p grpc-proto-artifact && \
protoc --go_out=./grpc-proto-artifact \
--proto_path=$(go env GOMODCACHE) \
--proto_path=$(go env GOPATH) \
--proto_path=. \
--go-grpc_out=./grpc-proto-artifact \
grpc-proto/grpc/reflection/v1/reflection.proto