package interceptor

import (
	"context"
	"fmt"
	"log"
	"strings"

	reflection "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	protobuf "google.golang.org/protobuf/proto"
	descriptor "google.golang.org/protobuf/types/descriptorpb"

	"github.com/kostyaBroTutor/auth_interceptor/pkg/contexts"
	"github.com/kostyaBroTutor/auth_interceptor/proto"
)

// getAuthenticatedMethodsForService returns a map with methods
// that require authentication.
// fullServiceName must be in the format "package.name.ServiceName".
func getAuthenticatedMethodsForService(
	ctx context.Context,
	fullServiceName string,
	reflectionClient reflection.ServerReflectionClient,
) (map[string]*proto.MethodAuthOptions, error) {
	reflectionInfo, err := reflectionClient.ServerReflectionInfo(
		contexts.ToOutgoing(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"can not to load server reflection info, error: %w", err,
		)
	}

	defer func() {
		if errr := reflectionInfo.CloseSend(); errr != nil {
			log.Printf("can not to close reflection info, error: %s", err)
		}
	}()

	service, err := findServicesDefinition(fullServiceName, reflectionInfo)
	if err != nil {
		return nil, fmt.Errorf(
			"can not to find services definition, error: %w", err,
		)
	}

	authenticatedMethods := make(map[string]*proto.MethodAuthOptions)

	for _, method := range service.GetMethod() {
		options := authenticationMethodsOptions(method.GetOptions())
		if options == nil {
			continue
		}

		fullRPCMethodName := toFullRPCMethodName(fullServiceName, method.GetName())
		authenticatedMethods[fullRPCMethodName] = options
	}

	return authenticatedMethods, nil
}

func findServicesDefinition(
	fullServiceName string,
	reflectionClient reflection.ServerReflection_ServerReflectionInfoClient,
) (*descriptor.ServiceDescriptorProto, error) {
	if err := reflectionClient.Send(&reflection.ServerReflectionRequest{
		MessageRequest: &reflection.ServerReflectionRequest_FileContainingSymbol{
			FileContainingSymbol: fullServiceName,
		},
	}); err != nil {
		return nil, fmt.Errorf(
			"error while sending request for reflection client: %w", err,
		)
	}

	response, err := reflectionClient.Recv()
	if err != nil {
		return nil, fmt.Errorf(
			"error while receive response from the reflection client: %w", err,
		)
	}

	packageName, shortServiceName := splitPackageAndServiceName(fullServiceName)

	for _, serializedProto := range response.
		GetFileDescriptorResponse().GetFileDescriptorProto() {
		var protoFileDescriptor descriptor.FileDescriptorProto
		if err = protobuf.Unmarshal(serializedProto, &protoFileDescriptor); err != nil {
			return nil, fmt.Errorf(
				"unmarhal file descriptor proto end with error: %w", err,
			)
		}

		// Find serviceDescriptor.
		for _, service := range protoFileDescriptor.GetService() {
			// shortPackageName does not contain the full package path,
			// so for reliability, we check
			// the full type names of arguments of all methods.
			if service.GetName() == shortServiceName &&
				isAllMethodsFromPackage(service, packageName) {
				return service, nil
			}
		}
	}

	log.Panic("service not found")

	return new(descriptor.ServiceDescriptorProto), nil
}

func splitPackageAndServiceName(
	fullServiceName string,
) (packageName string, serviceName string) {
	parts := strings.Split(fullServiceName, ".")

	return strings.Join(parts[:len(parts)-1], "."), parts[len(parts)-1]
}

func isAllMethodsFromPackage(
	service *descriptor.ServiceDescriptorProto, packageName string,
) bool {
	packageName = fmt.Sprintf(".%s.", packageName)

	for _, method := range service.GetMethod() {
		if !strings.HasPrefix(method.GetInputType(), packageName) {
			return false
		}
	}

	return true
}

func authenticationMethodsOptions(
	methodOptions protobuf.Message,
) *proto.MethodAuthOptions {
	extension := protobuf.GetExtension(methodOptions, proto.E_Options)

	options, isAuthOptions := extension.(*proto.MethodAuthOptions)
	if !isAuthOptions || !options.GetAuthenticated() {
		return nil
	}

	return options
}

// toFullServiceName provides full RPC method name
// in format /package.service/method.
func toFullRPCMethodName(
	fullServiceName string, shortMethodName string,
) string {
	return fmt.Sprintf("/%s/%s", fullServiceName, shortMethodName)
}
