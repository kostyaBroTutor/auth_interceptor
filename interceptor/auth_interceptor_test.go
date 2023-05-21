package interceptor_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"github.com/golang/protobuf/descriptor" //nolint:staticcheck // #TODO: use the google.golang.org/protobuf/reflect/protoreflect
	reflection "github.com/kostyaBroTutor/auth_interceptor/grpc-proto-artifact/google.golang.org/grpc/reflection/grpc_reflection_v1"
	"github.com/kostyaBroTutor/auth_interceptor/interceptor"
	mocks "github.com/kostyaBroTutor/auth_interceptor/interceptor/mocks"
	test "github.com/kostyaBroTutor/auth_interceptor/interceptor/testing"
	testing_mocks "github.com/kostyaBroTutor/auth_interceptor/interceptor/testing/mocks"
	"github.com/kostyaBroTutor/auth_interceptor/pkg/contexts"
	"github.com/kostyaBroTutor/auth_interceptor/proto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"testing"
)

//nolint:lll
//go:generate mockery --srcpkg=github.com/kostyaBroTutor/auth_interceptor/grpc-proto-artifact/google.golang.org/grpc/reflection/grpc_reflection_v1 --name=ServerReflectionClient
//go:generate mockery --srcpkg=github.com/kostyaBroTutor/auth_interceptor/grpc-proto-artifact/google.golang.org/grpc/reflection/grpc_reflection_v1 --name=ServerReflection_ServerReflectionInfoClient

const (
	testServiceName     = "TestService"
	testFullServiceName = "interceptor.test." + testServiceName

	testErrorMessage = "testError"
	testUserToken    = "testUserToken"
)

//nolint:gochecknoglobals
var (
	testContext   = metadata.NewIncomingContext(context.Background(), nil)
	testUserID    = "testUserID"
	testRequest   = new(interface{})
	testResponse  = new(interface{})
	testError     = errors.New(testErrorMessage) //nolint:golint,revive,errname,stylecheck
	testTokenInfo = &interceptor.TokenInfo{
		UserID: testUserID,
	}
	testTokenInfoWithPermissions = func(
		permissions ...proto.Permission,
	) *interceptor.TokenInfo {
		return &interceptor.TokenInfo{
			UserID:      testUserID,
			Permissions: permissions,
		}
	}
	testTokenInfoWithRoles = func(
		roles ...proto.Role,
	) *interceptor.TokenInfo {
		return &interceptor.TokenInfo{
			UserID: testUserID,
			Roles:  roles,
		}
	}
	testTokenInfoWithRolesAndPermissions = func(
		roles []proto.Role, permissions []proto.Permission,
	) *interceptor.TokenInfo {
		return &interceptor.TokenInfo{
			UserID:      testUserID,
			Roles:       roles,
			Permissions: permissions,
		}
	}

	testAuthenticatedMethodWithPermissionsOnlyCallInfo = &grpc.UnaryServerInfo{
		FullMethod: "/" + testFullServiceName + "/TestAuthenticatedMethodWithPermissionsOnly",
	}
	testAuthenticatedMethodWithRolesOnlyCallInfo = &grpc.UnaryServerInfo{
		FullMethod: "/" + testFullServiceName + "/TestAuthenticatedMethodWithRolesOnly",
	}
	testAuthenticatedMethodsWithRolesAndPermissionsCallInfo = &grpc.UnaryServerInfo{
		FullMethod: "/" + testFullServiceName + "/TestAuthenticatedMethodsWithRolesAndPermissions",
	}
	testAuthenticatedMethodNoPermissionsAndRolesCallInfo = &grpc.UnaryServerInfo{
		FullMethod: "/" + testFullServiceName + "/TestAuthenticatedMethodNoPermissionsAndRoles",
	}
	testNotAuthenticatedMethodExplicitCallInfo = &grpc.UnaryServerInfo{
		FullMethod: "/" + testFullServiceName + "/TestNotAuthenticatedMethodExplicit",
	}
	testNotAuthenticatedMethodNoAnnotationCallInfo = &grpc.UnaryServerInfo{
		FullMethod: "/" + testFullServiceName + "/TestNotAuthenticatedMethodNoAnnotation",
	}
)

func TestAuthInterceptor(t *testing.T) {
	t.Parallel()

	RegisterFailHandler(Fail)
	RunSpecs(t, "testing auth interceptor")
}

var _ = Describe("AuthInterceptor", func() {
	var (
		authClientMock              *mocks.AuthClient
		reflectionClientMock        *mocks.ServerReflectionClient
		handlerMock                 *testing_mocks.UnaryHandler
		authUnaryServiceInterceptor grpc.UnaryServerInterceptor
	)

	BeforeEach(func() {
		authClientMock = mocks.NewAuthClient(GinkgoT())
		reflectionClientMock = mocks.NewServerReflectionClient(GinkgoT())
		handlerMock = testing_mocks.NewUnaryHandler(GinkgoT())

		authUnaryServiceInterceptor = interceptor.NewAuthUnaryServiceInterceptor(
			authClientMock,
			testFullServiceName,
			reflectionClientMock,
		)
	})

	Describe("method with permissions only", func() {
		It("should return error if user is not authenticated", func() {
			testUserNotAuthenticated(
				reflectionClientMock, handlerMock, authUnaryServiceInterceptor,
			)
		})

		It("should return error if user does not have required permissions", func() {
			setupForAuthenticatedUserWithPermissionsAndRoles(
				authClientMock, reflectionClientMock, handlerMock,
				[]proto.Permission{
					proto.Permission_READ_SOMETHING_PERMISSION,
					proto.Permission_WRITE_SOMETHING_PERMISSION,
				},
				[]proto.Role{},
			)

			response, err := authUnaryServiceInterceptor(
				testContext, testRequest,
				testAuthenticatedMethodWithPermissionsOnlyCallInfo,
				handlerMock.GrpcUnaryHandler,
			)
			Expect(err).Should(HaveOccurred())
			Expect(response).To(BeNil())
			verifyError(
				err, codes.PermissionDenied,
				"Required permissions [permission:READ_SOMETHING_PERMISSION "+
					"permission:CHANGE_SOMETHING_PERMISSION] are not granted",
				"granted [permission:READ_SOMETHING_PERMISSION "+
					"permission:WRITE_SOMETHING_PERMISSION]",
			)
		})

		It("should not return error if user is authenticated and has required permissions", func() {
			setupForAuthenticatedUserWithPermissionsAndRoles(
				authClientMock, reflectionClientMock, handlerMock,
				[]proto.Permission{
					proto.Permission_READ_SOMETHING_PERMISSION,
					proto.Permission_CHANGE_SOMETHING_PERMISSION,
				},
				[]proto.Role{},
			)

			response, err := authUnaryServiceInterceptor(
				testContext, testRequest,
				testAuthenticatedMethodWithPermissionsOnlyCallInfo,
				handlerMock.GrpcUnaryHandler,
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(response).To(Equal(testResponse))
		})
	})

	Describe("method with roles only", func() {
		It("should return error if user is not authenticated", func() {
			testUserNotAuthenticated(
				reflectionClientMock, handlerMock, authUnaryServiceInterceptor,
			)
		})

		It("should return error if user does not have required roles", func() {
			setupForAuthenticatedUserWithPermissionsAndRoles(
				authClientMock, reflectionClientMock, handlerMock,
				[]proto.Permission{},
				[]proto.Role{proto.Role_EMPLOYEE_ROLE},
			)

			response, err := authUnaryServiceInterceptor(
				testContext, testRequest,
				testAuthenticatedMethodWithRolesOnlyCallInfo,
				handlerMock.GrpcUnaryHandler,
			)
			Expect(err).Should(HaveOccurred())
			Expect(response).To(BeNil())
			verifyError(
				err, codes.PermissionDenied,
				"Required roles [role:ADMIN_ROLE] are not granted",
				"granted [role:EMPLOYEE_ROLE]",
			)
		})

		It("should not return error if user is authenticated and has required roles", func() {
			setupForAuthenticatedUserWithPermissionsAndRoles(
				authClientMock, reflectionClientMock, handlerMock,
				[]proto.Permission{},
				[]proto.Role{proto.Role_ADMIN_ROLE},
			)

			response, err := authUnaryServiceInterceptor(
				testContext, testRequest,
				testAuthenticatedMethodWithRolesOnlyCallInfo,
				handlerMock.GrpcUnaryHandler,
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(response).ToNot(BeNil())
		})
	})

	Describe("method with roles and permissions", func() {
		It("should return error if user is not authenticated", func() {
			testUserNotAuthenticated(
				reflectionClientMock, handlerMock, authUnaryServiceInterceptor,
			)
		})

		It("should return error if user does not have required permissions", func() {
			setupForAuthenticatedUserWithPermissionsAndRoles(
				authClientMock, reflectionClientMock, handlerMock,
				[]proto.Permission{proto.Permission_READ_SOMETHING_PERMISSION},
				[]proto.Role{proto.Role_EMPLOYEE_ROLE},
			)

			response, err := authUnaryServiceInterceptor(
				testContext, testRequest,
				testAuthenticatedMethodsWithRolesAndPermissionsCallInfo,
				handlerMock.GrpcUnaryHandler,
			)
			Expect(response).To(BeNil())
			Ω(err).Should(HaveOccurred())
			verifyError(
				err, codes.PermissionDenied,
				"Required permissions [permission:WRITE_SOMETHING_PERMISSION] are not granted",
				"granted [permission:READ_SOMETHING_PERMISSION]",
			)
		})

		It("should return error if user does not have required roles", func() {
			setupForAuthenticatedUserWithPermissionsAndRoles(
				authClientMock, reflectionClientMock, handlerMock,
				[]proto.Permission{proto.Permission_WRITE_SOMETHING_PERMISSION},
				[]proto.Role{proto.Role_ADMIN_ROLE},
			)

			response, err := authUnaryServiceInterceptor(
				testContext, testRequest,
				testAuthenticatedMethodsWithRolesAndPermissionsCallInfo,
				handlerMock.GrpcUnaryHandler,
			)
			Expect(response).To(BeNil())
			Ω(err).Should(HaveOccurred())
			verifyError(
				err, codes.PermissionDenied,
				"Required roles [role:EMPLOYEE_ROLE] are not granted",
				"granted [role:ADMIN_ROLE]",
			)
		})

		It("should not return error if user is authenticated and has required permissions and roles", func() {
			setupForAuthenticatedUserWithPermissionsAndRoles(
				authClientMock, reflectionClientMock, handlerMock,
				[]proto.Permission{proto.Permission_WRITE_SOMETHING_PERMISSION},
				[]proto.Role{proto.Role_EMPLOYEE_ROLE},
			)

			response, err := authUnaryServiceInterceptor(
				testContext, testRequest,
				testAuthenticatedMethodsWithRolesAndPermissionsCallInfo,
				handlerMock.GrpcUnaryHandler,
			)
			Expect(err).To(BeNil())
			Expect(response).To(Equal(testResponse))
		})
	})

	Describe("method without required permissions and roles", func() {
		It("should return error if user is not authenticated", func() {
			testUserNotAuthenticated(
				reflectionClientMock, handlerMock, authUnaryServiceInterceptor,
			)
		})

		It("should not return error if user is authenticated", func() {
			setupForAuthenticatedUserWithoutPermissionsAndRoles(
				authClientMock, reflectionClientMock, handlerMock,
			)

			response, err := authUnaryServiceInterceptor(
				testContext, testRequest,
				testAuthenticatedMethodNoPermissionsAndRolesCallInfo,
				handlerMock.GrpcUnaryHandler,
			)
			Expect(response).To(Equal(testResponse))
			Expect(err).To(Equal(testError))
		})
	})

	Describe("method with explicit not authenticated annotation", func() {
		It("should not return error if user is not authenticated", func() {
			testUserNotAuthenticated(
				reflectionClientMock, handlerMock, authUnaryServiceInterceptor,
			)
		})

		It("should not return error if user is authenticated", func() {
			setupForAuthenticatedUserWithoutPermissionsAndRoles(
				authClientMock, reflectionClientMock, handlerMock,
			)

			response, err := authUnaryServiceInterceptor(
				testContext, testRequest,
				testNotAuthenticatedMethodExplicitCallInfo,
				handlerMock.GrpcUnaryHandler,
			)

			Expect(response).To(Equal(testResponse))
			Expect(err).To(Equal(testError))
		})
	})

	Describe("method without annotation", func() {
		It("should not return error if user is not authenticated", func() {
			testUserNotAuthenticated(
				reflectionClientMock, handlerMock, authUnaryServiceInterceptor,
			)
		})

		It("should not return error if user is authenticated", func() {
			setupForAuthenticatedUserWithoutPermissionsAndRoles(
				authClientMock, reflectionClientMock, handlerMock,
			)
			testContext = metadata.NewIncomingContext(
				testContext,
				metadata.Pairs(
					string(interceptor.AuthTokenMetadataName),
					testUserToken,
				),
			)

			response, err := authUnaryServiceInterceptor(
				testContext, testRequest,
				testNotAuthenticatedMethodNoAnnotationCallInfo,
				handlerMock.GrpcUnaryHandler,
			)

			Expect(response).To(Equal(testResponse))
			Expect(err).To(Equal(testError))
		})
	})

	Describe("cache", func() {
		It("allow not authenticated request using cache", func() {
			makeCallToCreateCache(
				authUnaryServiceInterceptor, reflectionClientMock, handlerMock,
			)

			handlerMock.On("GrpcUnaryHandler", testContext, testRequest).
				Return(testResponse, testError)

			response, err := authUnaryServiceInterceptor(
				testContext, testRequest,
				testNotAuthenticatedMethodNoAnnotationCallInfo,
				handlerMock.GrpcUnaryHandler,
			)
			Expect(response).To(Equal(testResponse))
			Expect(err).To(Equal(testError))
		})

		It("allow authenticated request using cache", func() {
			makeCallToCreateCache(
				authUnaryServiceInterceptor, reflectionClientMock, handlerMock,
			)

			setupForAuthenticatedUserWithPermissionsAndRoles(
				authClientMock, reflectionClientMock, handlerMock,
				[]proto.Permission{proto.Permission_WRITE_SOMETHING_PERMISSION},
				[]proto.Role{proto.Role_EMPLOYEE_ROLE},
			)

			response, err := authUnaryServiceInterceptor(
				testContext, testRequest,
				testAuthenticatedMethodsWithRolesAndPermissionsCallInfo,
				handlerMock.GrpcUnaryHandler,
			)
			Expect(err).To(BeNil())
			Expect(response).To(Equal(testResponse))
		})
	})

	Describe("error cases", func() {
		It("provide error if metadata invalid", func() {
			makeCallToCreateCache(
				authUnaryServiceInterceptor, reflectionClientMock, handlerMock,
			)

			_, err := authUnaryServiceInterceptor(
				context.Background(), testRequest,
				testAuthenticatedMethodsWithRolesAndPermissionsCallInfo,
				handlerMock.GrpcUnaryHandler,
			)
			verifyError(err, codes.Unauthenticated, "invalid request metadata")
		})

		It("provide error if missing token", func() {
			testContextWithMissingToken := metadata.NewIncomingContext(
				testContext,
				metadata.Pairs("testKey", "testValue"),
			)
			setupReflection(testContextWithMissingToken, reflectionClientMock)

			_, err := authUnaryServiceInterceptor(
				testContextWithMissingToken, testRequest,
				testAuthenticatedMethodsWithRolesAndPermissionsCallInfo,
				handlerMock.GrpcUnaryHandler,
			)
			verifyError(err, codes.Unauthenticated, "missing auth token")
		})

		It("provide error if token invalid", func() {
			testContextInvalidToken := metadata.NewIncomingContext(
				testContext,
				metadata.Pairs(
					string(interceptor.AuthTokenMetadataName), "testValue1",
					string(interceptor.AuthTokenMetadataName), "testValue2",
				),
			)
			setupReflection(testContextInvalidToken, reflectionClientMock)
			_, err := authUnaryServiceInterceptor(
				testContextInvalidToken, testRequest,
				testAuthenticatedMethodsWithRolesAndPermissionsCallInfo,
				handlerMock.GrpcUnaryHandler,
			)
			verifyError(err, codes.Unauthenticated, "invalid auth token")
		})

		It("provide error if auth with token was failed", func() {
			testContextWithAuthToken := metadata.NewIncomingContext(
				testContext,
				metadata.Pairs(
					string(interceptor.AuthTokenMetadataName),
					testUserToken,
				),
			)
			setupReflection(testContextWithAuthToken, reflectionClientMock)
			authClientMock.On("Auth",
				contexts.ToOutgoing(testContextWithAuthToken),
				testUserToken,
			).Return(nil, testError)

			_, err := authUnaryServiceInterceptor(
				testContextWithAuthToken, testRequest,
				testAuthenticatedMethodsWithRolesAndPermissionsCallInfo,
				handlerMock.GrpcUnaryHandler,
			)
			verifyError(err, codes.Unauthenticated, testErrorMessage)
		})
	})
})

func testUserNotAuthenticated(
	reflectionClientMock *mocks.ServerReflectionClient,
	handlerMock *testing_mocks.UnaryHandler,
	authUnaryServiceInterceptor grpc.UnaryServerInterceptor,
) {
	setupForNonAuthenticatedUser(reflectionClientMock, handlerMock)

	response, err := authUnaryServiceInterceptor(
		testContext, testRequest,
		testAuthenticatedMethodsWithRolesAndPermissionsCallInfo,
		handlerMock.GrpcUnaryHandler,
	)
	Expect(response).To(BeNil())
	Ω(err).Should(HaveOccurred())
	verifyError(err, codes.Unauthenticated, "user is not authenticated")
}

func setupForNonAuthenticatedUser(
	reflectionClientMock *mocks.ServerReflectionClient,
	handlerMock *testing_mocks.UnaryHandler,
) {
	setupReflection(testContext, reflectionClientMock)
	handlerMock.On("GrpcUnaryHandler", testContext, testRequest).
		Return(testResponse, testError)
}

func setupForAuthenticatedUserWithoutPermissionsAndRoles(
	authClientMock *mocks.AuthClient,
	reflectionClientMock *mocks.ServerReflectionClient,
	handlerMock *testing_mocks.UnaryHandler,
) {
	testAuthContext := metadata.NewIncomingContext(
		testContext,
		metadata.Pairs(
			string(interceptor.AuthTokenMetadataName), testUserToken,
		),
	)
	setupAuthResponse(
		contexts.ToOutgoing(testAuthContext), authClientMock, testTokenInfo,
	)

	handlerContext := context.WithValue(
		testAuthContext, interceptor.UserIDMetadataName, testUserID,
	)
	handlerContext = context.WithValue(
		handlerContext, interceptor.RolesMetadataName, []proto.Role{},
	)
	handlerContext = context.WithValue(
		handlerContext, interceptor.PermissionMetadataName, []proto.Permission{},
	)

	setupReflection(testAuthContext, reflectionClientMock)
	handlerMock.On("GrpcUnaryHandler", handlerContext, testRequest).
		Return(testResponse, testError)
}

func setupForAuthenticatedUserWithPermissionsAndRoles(
	authClientMock *mocks.AuthClient,
	reflectionClientMock *mocks.ServerReflectionClient,
	handlerMock *testing_mocks.UnaryHandler,
	permissions []proto.Permission,
	roles []proto.Role,
) {
	testAuthContext := metadata.NewIncomingContext(
		testContext,
		metadata.Pairs(
			string(interceptor.AuthTokenMetadataName), testUserToken,
		),
	)
	setupAuthResponse(
		contexts.ToOutgoing(testAuthContext), authClientMock, testTokenInfo,
	)

	handlerContext := context.WithValue(
		testAuthContext, interceptor.UserIDMetadataName, testUserID,
	)
	handlerContext = context.WithValue(
		handlerContext, interceptor.RolesMetadataName, roles,
	)
	handlerContext = context.WithValue(
		handlerContext, interceptor.PermissionMetadataName, permissions,
	)

	setupReflection(testAuthContext, reflectionClientMock)
	handlerMock.On("GrpcUnaryHandler", handlerContext, testRequest).
		Return(testResponse, testError)
}

// setupReflection service to return
// the serialized file descriptor with a definition of the testFullServiceName.
func setupReflection(
	ctx context.Context,
	reflectionClientMock *mocks.ServerReflectionClient,
) {
	message := test.TestMessage{}
	fileDescriptorGziped, _ := descriptor.MessageRawDescriptor(
		message.ProtoReflect().Interface())
	zr, err := gzip.NewReader(bytes.NewBuffer(fileDescriptorGziped))
	Ω(err).ShouldNot(HaveOccurred())

	fileDescriptorSerialized, err := io.ReadAll(zr)
	Ω(err).ShouldNot(HaveOccurred())
	Ω(zr.Close()).ShouldNot(HaveOccurred())

	protoStream := new(mocks.ServerReflection_ServerReflectionInfoClient)
	protoStream.On(
		"Send",
		&reflection.ServerReflectionRequest{
			MessageRequest: &reflection.ServerReflectionRequest_FileContainingSymbol{
				FileContainingSymbol: testFullServiceName,
			},
		},
	).Return(nil)
	protoStream.On(
		"Recv",
	).Return(&reflection.ServerReflectionResponse{
		MessageResponse: &reflection.ServerReflectionResponse_FileDescriptorResponse{
			FileDescriptorResponse: &reflection.FileDescriptorResponse{
				FileDescriptorProto: [][]byte{fileDescriptorSerialized},
			},
		}},
		nil,
	)
	protoStream.On("CloseSend").Return(nil)

	reflectionClientMock.On(
		"ServerReflectionInfo",
		contexts.ToOutgoing(ctx),
	).Return(protoStream, nil)
}

func setupAuthResponse(
	ctx context.Context,
	mockAuthClient *mocks.AuthClient,
	info *interceptor.TokenInfo,
) {
	mockAuthClient.On("Auth", ctx, testUserToken).Return(info, nil)
}

func makeCallToCreateCache(
	authInterceptor grpc.UnaryServerInterceptor,
	reflectionClientMock *mocks.ServerReflectionClient,
	callEndHandlerMock *testing_mocks.UnaryHandler,
) {
	setupReflection(testContext, reflectionClientMock)
	callEndHandlerMock.On("GrpcUnaryHandler", testContext, testRequest).
		Return(testResponse, testError)

	response, err := authInterceptor(
		testContext, testRequest,
		testNotAuthenticatedMethodNoAnnotationCallInfo,
		callEndHandlerMock.GrpcUnaryHandler,
	)
	Expect(response).To(Equal(testResponse))
	Expect(err).To(Equal(testError))
}

func verifyError(err error, code codes.Code, messages ...string) {
	grpcErr, ok := status.FromError(err)
	Expect(ok).To(BeTrue())
	Expect(grpcErr.Code()).To(Equal(code))

	for _, message := range messages {
		Expect(grpcErr.Message()).To(ContainSubstring(message))
	}
}
