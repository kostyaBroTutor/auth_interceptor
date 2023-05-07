#!/bin/sh
echo "running ./proto/build.sh"
./proto/build.sh
echo "running go generate ./..."
go generate ./...
