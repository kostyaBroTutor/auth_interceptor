#!/bin/sh
echo "running build_grpc-proto.sh"
./scripts/build_grpc-proto.sh
echo "running ./proto/build.sh"
./proto/build.sh
echo "running ./example/service/proto/build.sh"
./example/service/proto/build.sh
echo "running ./interceptor/testdata/build.sh"
./interceptor/testing/build.sh
echo "running go generate ./..."
go generate ./...
