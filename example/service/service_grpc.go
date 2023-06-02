package service

import (
	"fmt"
	"log"
	"net"
	"time"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	proto "github.com/kostyaBroTutor/auth_interceptor/example/service/proto"
	grpcreflectionv1 "github.com/kostyaBroTutor/auth_interceptor/grpc-proto-artifact/google.golang.org/grpc/reflection/grpc_reflection_v1"
	"github.com/kostyaBroTutor/auth_interceptor/interceptor"
	"github.com/kostyaBroTutor/auth_interceptor/pkg/process"
)

type Config struct {
	ListenAddr           string
	AuthClient           interceptor.AuthClient
	ExampleServiceServer proto.ExampleServiceServer
}

// NewExampleGrpcServer starts new gRPC server with ExampleService.
// IMPORTANT: It returns closer that should be called to stop server.
//
//nolint:funlen
func NewExampleGrpcServer(config Config) (func(), error) {
	listenerForExampleServer, err := net.Listen("tcp", config.ListenAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	exampleServiceClientConn := new(grpc.ClientConn)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmiddleware.ChainUnaryServer(
				interceptor.NewAuthUnaryServiceInterceptor(
					config.AuthClient,
					"example.proto.ExampleService",
					grpcreflectionv1.NewServerReflectionClient(
						exampleServiceClientConn,
					),
				),
			),
		),
		grpc.Creds(insecure.NewCredentials()),
	)
	proto.RegisterExampleServiceServer(grpcServer, config.ExampleServiceServer)
	reflection.Register(grpcServer)

	closer := func() {
		if err := exampleServiceClientConn.Close(); err != nil {
			log.Println("failed to close example service client connection, error: " + err.Error())
		}

		grpcServer.GracefulStop()
	}

	go func() {
		if err := grpcServer.Serve(listenerForExampleServer); err != nil {
			log.Println("failed to serve example server, error: " + err.Error())

			process.Terminate()
		}
	}()

	time.Sleep(1 * time.Second)

	exampleServiceClientConnNew, err := grpc.Dial(
		config.ListenAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to connect to odo acceptance service: %w", err,
		)
	}

	*exampleServiceClientConn = *exampleServiceClientConnNew //nolint:govet

	return closer, nil
}
