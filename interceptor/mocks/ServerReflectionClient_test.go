// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"

	grpc_reflection_v1alpha "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"

	mock "github.com/stretchr/testify/mock"
)

// ServerReflectionClient is an autogenerated mock type for the ServerReflectionClient type
type ServerReflectionClient struct {
	mock.Mock
}

type ServerReflectionClient_Expecter struct {
	mock *mock.Mock
}

func (_m *ServerReflectionClient) EXPECT() *ServerReflectionClient_Expecter {
	return &ServerReflectionClient_Expecter{mock: &_m.Mock}
}

// ServerReflectionInfo provides a mock function with given fields: ctx, opts
func (_m *ServerReflectionClient) ServerReflectionInfo(ctx context.Context, opts ...grpc.CallOption) (grpc_reflection_v1alpha.ServerReflection_ServerReflectionInfoClient, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 grpc_reflection_v1alpha.ServerReflection_ServerReflectionInfoClient
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...grpc.CallOption) (grpc_reflection_v1alpha.ServerReflection_ServerReflectionInfoClient, error)); ok {
		return rf(ctx, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...grpc.CallOption) grpc_reflection_v1alpha.ServerReflection_ServerReflectionInfoClient); ok {
		r0 = rf(ctx, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(grpc_reflection_v1alpha.ServerReflection_ServerReflectionInfoClient)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ServerReflectionClient_ServerReflectionInfo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ServerReflectionInfo'
type ServerReflectionClient_ServerReflectionInfo_Call struct {
	*mock.Call
}

// ServerReflectionInfo is a helper method to define mock.On call
//   - ctx context.Context
//   - opts ...grpc.CallOption
func (_e *ServerReflectionClient_Expecter) ServerReflectionInfo(ctx interface{}, opts ...interface{}) *ServerReflectionClient_ServerReflectionInfo_Call {
	return &ServerReflectionClient_ServerReflectionInfo_Call{Call: _e.mock.On("ServerReflectionInfo",
		append([]interface{}{ctx}, opts...)...)}
}

func (_c *ServerReflectionClient_ServerReflectionInfo_Call) Run(run func(ctx context.Context, opts ...grpc.CallOption)) *ServerReflectionClient_ServerReflectionInfo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *ServerReflectionClient_ServerReflectionInfo_Call) Return(_a0 grpc_reflection_v1alpha.ServerReflection_ServerReflectionInfoClient, _a1 error) *ServerReflectionClient_ServerReflectionInfo_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ServerReflectionClient_ServerReflectionInfo_Call) RunAndReturn(run func(context.Context, ...grpc.CallOption) (grpc_reflection_v1alpha.ServerReflection_ServerReflectionInfoClient, error)) *ServerReflectionClient_ServerReflectionInfo_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewServerReflectionClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewServerReflectionClient creates a new instance of ServerReflectionClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewServerReflectionClient(t mockConstructorTestingTNewServerReflectionClient) *ServerReflectionClient {
	mock := &ServerReflectionClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
