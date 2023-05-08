#!/bin/sh
echo "running ./proto/build.sh"
./proto/build.sh
echo "running ./interceptor/testdata/build.sh"
./interceptor/testing/build.sh
echo "running go generate ./..."
go generate ./...
