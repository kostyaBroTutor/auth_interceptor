#!/bin/sh
protoc --go_out=. \
--proto_path=$(go env GOMODCACHE) \
--proto_path=/workspace/protobuf/src \
--proto_path=. \
--go-grpc_out=. \
proto/annotation.proto
