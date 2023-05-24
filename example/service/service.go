package service

import (
	"context"
	"log"

	proto "github.com/kostyaBroTutor/auth_interceptor/example/service/proto"
)

type ExampleService struct {
	proto.UnimplementedExampleServiceServer
}

func NewExampleService() *ExampleService {
	return new(ExampleService)
}

func (e ExampleService) FreeMethod(
	context.Context, *proto.TestMessage,
) (*proto.TestMessage, error) {
	log.Println("FreeMethod was called")

	return new(proto.TestMessage), nil
}

func (e ExampleService) LimitAccessMethod(
	context.Context, *proto.TestMessage,
) (*proto.TestMessage, error) {
	log.Println("LimitAccessMethod was called")

	return new(proto.TestMessage), nil
}
