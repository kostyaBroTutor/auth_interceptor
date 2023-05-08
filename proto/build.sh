#!/bin/sh
protoc --go_out=. \
--proto_path=$(go env GOMODCACHE) \
--proto_path=/workspace/protobuf/src \
--proto_path=. \
proto/annotation.proto
