// Package contexts contains method for working with context.
package contexts

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// ToOutgoing converts the incoming context to the outgoing context
// with saving metadata.
func ToOutgoing(srcCtx context.Context) context.Context {
	requestMetadata, exists := metadata.FromIncomingContext(srcCtx)
	if !exists {
		panic("Invalid incoming context")
	}

	return metadata.NewOutgoingContext(srcCtx, requestMetadata)
}
