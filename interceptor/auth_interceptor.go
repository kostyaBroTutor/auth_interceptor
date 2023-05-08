// Package interceptor provides the middleware
// for check roles and permissions of user.
package interceptor

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	reflection "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"

	"github.com/kostyaBroTutor/auth_interceptor/pkg/contexts"
	"github.com/kostyaBroTutor/auth_interceptor/proto"
)

//go:generate mockery --name=AuthClient

type AuthClient interface {
	Auth(ctx context.Context, token string) (*TokenInfo, error)
}

type MetadataName string

const (
	AuthTokenMetadataName  = MetadataName("auth-token")
	UserIDMetadataName     = MetadataName("user-id")
	RolesMetadataName      = MetadataName("roles")
	PermissionMetadataName = MetadataName("permissions")
)

type authInterceptor struct {
	authClient       AuthClient
	fullServiceName  string
	reflectionClient reflection.ServerReflectionClient

	authenticatedMethodsMutex sync.Mutex
	authenticatedMethods      map[string]*proto.MethodAuthOptions
}

func NewAuthUnaryServiceInterceptor(
	authClient AuthClient,
	// fullServiceName is the package name from the proto file plus service name.
	fullServiceName string,
	reflectionClient reflection.ServerReflectionClient,
) grpc.UnaryServerInterceptor {
	interceptorObject := &authInterceptor{
		authClient:       authClient,
		fullServiceName:  fullServiceName,
		reflectionClient: reflectionClient,
		// No need to init map, lazy initialization in ai.getAuthenticatedMethods()
	}

	return interceptorObject.AuthUnaryServiceInterceptor
}

// AuthUnaryServiceInterceptor checks the token specified in the metadata.
// If the token is valid, userID and the payload will be added to the metadata
// for the following query handlers.
// The login request does not require a token and
// will be forwarded to the following handlers without verification.
func (ai *authInterceptor) AuthUnaryServiceInterceptor(
	requestContext context.Context,
	request interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	authenticatedMethods, err := ai.getAuthenticatedMethods(requestContext)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			fmt.Sprintf("can not to get autenticated methods, error: %s", err),
		)
	}

	authenticationOptions, needAuthentication :=
		authenticatedMethods[info.FullMethod]
	if !needAuthentication {
		return handler(requestContext, request)
	}

	authTokenInfo, err := ai.verifyToken(requestContext)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	if !authTokenInfo.HasAccessToMethodWithWithAuthOptions(authenticationOptions) {
		return nil, status.Error(codes.PermissionDenied, fmt.Sprintf(
			"Required permissions %+v, but only granted %+v",
			authenticationOptions, authTokenInfo,
		))
	}

	requestContext = ai.addTokenInfoToContext(requestContext, authTokenInfo)

	return handler(requestContext, request)
}

func (ai *authInterceptor) addTokenInfoToContext(
	ctx context.Context,
	tokenInfo *TokenInfo,
) context.Context {
	ctx = context.WithValue(ctx, UserIDMetadataName, tokenInfo.UserID)
	ctx = context.WithValue(ctx, RolesMetadataName, tokenInfo.Roles)
	ctx = context.WithValue(ctx, PermissionMetadataName, tokenInfo.Permissions)

	return ctx
}

func (ai *authInterceptor) verifyToken(
	requestContext context.Context,
) (*TokenInfo, error) {
	requestMetadata, exists := metadata.FromIncomingContext(requestContext)
	if !exists {
		return nil, errors.Join(ErrVerifyToken, ErrInvalidRequestMetadata)
	}

	authToken, exists := requestMetadata[string(AuthTokenMetadataName)]
	if !exists {
		return nil, errors.Join(ErrVerifyToken, ErrMissingAuthToken)
	}

	if len(authToken) != 1 {
		return nil, errors.Join(ErrVerifyToken, ErrInvalidAuthToken)
	}

	response, err := ai.authClient.Auth(
		contexts.ToOutgoing(requestContext), authToken[0],
	)
	if err != nil {
		return nil, errors.Join(ErrVerifyToken, err)
	}

	return response, nil
}

// getAuthenticatedMethods provide hash-table
// that describes the authentication requirements for service methods.
// For the first call, the reflection server
// will be used to get authentication information in methods.
// For the rest call, result will be return from the saved cache.
func (ai *authInterceptor) getAuthenticatedMethods(
	ctx context.Context,
) (map[string]*proto.MethodAuthOptions, error) {
	ai.authenticatedMethodsMutex.Lock()
	defer ai.authenticatedMethodsMutex.Unlock()

	if ai.authenticatedMethods != nil {
		return ai.authenticatedMethods, nil
	}

	authenticatedMethods, err := getAuthenticatedMethodsForService(
		ctx, ai.fullServiceName, ai.reflectionClient,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"can not to get authenticated methods for service %q, error: %w",
			ai.fullServiceName, err,
		)
	}

	ai.authenticatedMethods = authenticatedMethods

	return ai.authenticatedMethods, nil
}
