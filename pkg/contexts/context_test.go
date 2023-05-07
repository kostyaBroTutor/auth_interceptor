package contexts_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/metadata"

	"github.com/kostyaBroTutor/auth_interceptor/pkg/contexts"
)

func TestContexts(t *testing.T) {
	t.Parallel()

	RegisterFailHandler(Fail)
	RunSpecs(t, "testing contexts package")
}

var _ = It("can convert incoming context to outgoing context", func() {
	testMetadata := metadata.New(map[string]string{
		"test_key": "test_value",
	})
	testIncomingContext := metadata.NewIncomingContext(
		context.Background(), testMetadata,
	)
	expectedOutgoingContext := metadata.NewOutgoingContext(
		testIncomingContext, testMetadata,
	)

	Expect(contexts.ToOutgoing(testIncomingContext)).
		To(Equal(expectedOutgoingContext))
})

var _ = It("panic if context is wrong", func() {
	wrongContext := metadata.NewOutgoingContext(
		context.Background(), metadata.New(map[string]string{}),
	)

	Ω(func() { contexts.ToOutgoing(wrongContext) }).Should(Panic())
})
