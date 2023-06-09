// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	context "context"

	interceptor "github.com/kostyaBroTutor/auth_interceptor/interceptor"
	mock "github.com/stretchr/testify/mock"
)

// AuthClient is an autogenerated mock type for the AuthClient type
type AuthClient struct {
	mock.Mock
}

type AuthClient_Expecter struct {
	mock *mock.Mock
}

func (_m *AuthClient) EXPECT() *AuthClient_Expecter {
	return &AuthClient_Expecter{mock: &_m.Mock}
}

// Auth provides a mock function with given fields: ctx, token
func (_m *AuthClient) Auth(ctx context.Context, token string) (*interceptor.TokenInfo, error) {
	ret := _m.Called(ctx, token)

	var r0 *interceptor.TokenInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*interceptor.TokenInfo, error)); ok {
		return rf(ctx, token)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *interceptor.TokenInfo); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*interceptor.TokenInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthClient_Auth_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Auth'
type AuthClient_Auth_Call struct {
	*mock.Call
}

// Auth is a helper method to define mock.On call
//   - ctx context.Context
//   - token string
func (_e *AuthClient_Expecter) Auth(ctx interface{}, token interface{}) *AuthClient_Auth_Call {
	return &AuthClient_Auth_Call{Call: _e.mock.On("Auth", ctx, token)}
}

func (_c *AuthClient_Auth_Call) Run(run func(ctx context.Context, token string)) *AuthClient_Auth_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *AuthClient_Auth_Call) Return(_a0 *interceptor.TokenInfo, _a1 error) *AuthClient_Auth_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AuthClient_Auth_Call) RunAndReturn(run func(context.Context, string) (*interceptor.TokenInfo, error)) *AuthClient_Auth_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewAuthClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewAuthClient creates a new instance of AuthClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuthClient(t mockConstructorTestingTNewAuthClient) *AuthClient {
	mock := &AuthClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
