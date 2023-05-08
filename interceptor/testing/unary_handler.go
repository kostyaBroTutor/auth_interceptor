// Package testing contains mocks and service definitions for testing.
package testing

//go:generate mockery --name=UnaryHandler

import (
	"context"
)

// UnaryHandler wrap grpc.UnaryHandler function for testing.
type UnaryHandler interface {
	GrpcUnaryHandler(ctx context.Context, req interface{}) (interface{}, error)
}
