#!/bin/sh
protoc --go_out=. \
--proto_path=$(go env GOMODCACHE) \
--proto_path=$(go env GOPATH) \
--proto_path=/workspace/protobuf/src \
--proto_path=. \
--go-grpc_out=. \
--go_opt=Mproto/annotation.proto=github.com/kostyaBroTutor/auth_interceptor/proto \
example/service/proto/service.proto

